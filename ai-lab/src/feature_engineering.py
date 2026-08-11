import math
import logging
from datetime import datetime
from typing import Dict, List, Any, Tuple
import numpy as np

logger = logging.getLogger("feature_engineering")

class TimeSeriesFeatureBuilder:
    """
    Constructs time-series lag and rolling statistics features for demand forecasting.
    """
    def __init__(self, lags: List[int] = None, rolling_windows: List[int] = None):
        self.lags = lags or [1, 7, 14, 30]
        self.rolling_windows = rolling_windows or [7, 30]

    def build_features(self, raw_data: Dict[str, List[Dict[str, Any]]]) -> Tuple[List[str], np.ndarray, np.ndarray]:
        order_items = raw_data.get("order_items", [])
        orders = raw_data.get("orders", [])
        products = raw_data.get("products", [])

        # Build date to order lookup
        order_date_map = {}
        for o in orders:
            # Extract YYYY-MM-DD
            dt_str = o["created_at"].split(" ")[0]
            order_date_map[o["order_id"]] = dt_str

        product_map = {p["product_id"]: p for p in products}

        # Daily sales aggregation: (date, product_id) -> units_sold
        daily_sales: Dict[Tuple[str, str], float] = {}
        dates_set = set()
        product_ids_set = set(product_map.keys())

        for item in order_items:
            oid = item["order_id"]
            if oid in order_date_map:
                dt_str = order_date_map[oid]
                pid = item["product_id"]
                dates_set.add(dt_str)
                key = (dt_str, pid)
                daily_sales[key] = daily_sales.get(key, 0.0) + float(item["quantity"])

        sorted_dates = sorted(list(dates_set))
        if len(sorted_dates) < 35:
            logger.warning("Fewer than 35 days available for time-series feature engineering.")

        # Build daily matrix per product
        feature_rows = []
        target_rows = []

        feature_names = [
            "gross_margin_pct", "view_count", "day_of_week", "month", "is_weekend"
        ]
        for l in self.lags:
            feature_names.append(f"units_sold_lag_{l}")
        for w in self.rolling_windows:
            feature_names.append(f"rolling_{w}d_mean")
            feature_names.append(f"rolling_{w}d_std")

        date_dt_map = {d: datetime.strptime(d, "%Y-%m-%d") for d in sorted_dates}

        for idx, current_date_str in enumerate(sorted_dates):
            # Skip first 30 days so full lags can be populated
            if idx < 30 or idx >= len(sorted_dates) - 7:
                continue

            current_dt = date_dt_map[current_date_str]
            day_of_week = current_dt.weekday()
            month = current_dt.month
            is_weekend = 1.0 if day_of_week >= 5 else 0.0

            for pid in sorted(list(product_ids_set)):
                p_meta = product_map.get(pid, {})
                margin = p_meta.get("gross_margin_pct", 0.3)
                views = p_meta.get("view_count", 500)

                row_feats = [margin, float(views), float(day_of_week), float(month), is_weekend]

                # Lags
                for l in self.lags:
                    prev_date_str = sorted_dates[idx - l]
                    val = daily_sales.get((prev_date_str, pid), 0.0)
                    row_feats.append(val)

                # Rolling statistics
                for w in self.rolling_windows:
                    window_dates = sorted_dates[max(0, idx - w):idx]
                    vals = [daily_sales.get((d, pid), 0.0) for d in window_dates]
                    r_mean = float(np.mean(vals)) if vals else 0.0
                    r_std = float(np.std(vals)) if vals else 0.0
                    row_feats.append(r_mean)
                    row_feats.append(r_std)

                # Target horizon (next 7 days total units sold)
                future_dates = sorted_dates[idx + 1 : idx + 8]
                target_val = sum(daily_sales.get((fd, pid), 0.0) for fd in future_dates)

                feature_rows.append(row_feats)
                target_rows.append(target_val)

        X = np.array(feature_rows, dtype=np.float32) if feature_rows else np.empty((0, len(feature_names)), dtype=np.float32)
        y = np.array(target_rows, dtype=np.float32) if target_rows else np.empty((0,), dtype=np.float32)

        logger.info(f"Generated Demand Forecasting dataset: X.shape={X.shape}, y.shape={y.shape}")
        return feature_names, X, y


class InventoryStockoutFeatureBuilder:
    """
    Constructs stockout risk and reorder urgency features per product & shop.
    """
    def build_features(self, raw_data: Dict[str, List[Dict[str, Any]]]) -> Tuple[List[str], np.ndarray, np.ndarray]:
        inventory = raw_data.get("inventory", [])
        products = raw_data.get("products", [])
        order_items = raw_data.get("order_items", [])

        product_map = {p["product_id"]: p for p in products}

        # Calculate 7d sales burn rate per product
        total_qty_map: Dict[str, float] = {}
        for item in order_items:
            pid = item["product_id"]
            total_qty_map[pid] = total_qty_map.get(pid, 0.0) + float(item["quantity"])

        num_days = max(1, len(raw_data.get("orders", [])) // 15)

        feature_names = [
            "stock", "reserved_stock", "stock_burn_rate_7d",
            "supplier_lead_time_days", "estimated_days_to_stockout", "reorder_urgency_ratio"
        ]

        feature_rows = []
        target_rows = []

        for inv in inventory:
            pid = inv["product_id"]
            stock = float(inv["stock"])
            reserved = float(inv["reserved_stock"])

            p_meta = product_map.get(pid, {})
            lead_time = float(p_meta.get("supplier_lead_time_days", 3))

            daily_burn = (total_qty_map.get(pid, 0.0) / float(num_days)) / 3.0 # per shop approx
            est_days = stock / (daily_burn + 1e-5)
            urgency = lead_time / (est_days + 1e-5)

            row = [stock, reserved, daily_burn, lead_time, est_days, urgency]

            # Target binary flag: 1 if stockout within lead time, else 0
            target = 1.0 if est_days <= lead_time else 0.0

            feature_rows.append(row)
            target_rows.append(target)

        X = np.array(feature_rows, dtype=np.float32) if feature_rows else np.empty((0, len(feature_names)), dtype=np.float32)
        y = np.array(target_rows, dtype=np.float32) if target_rows else np.empty((0,), dtype=np.float32)

        logger.info(f"Generated Stockout Risk dataset: X.shape={X.shape}, y.shape={y.shape}")
        return feature_names, X, y


class OperationalAnomalyFeatureBuilder:
    """
    Constructs normalized tabular anomaly feature matrix from payment and fulfillment latency metrics.
    """
    def build_features(self, raw_data: Dict[str, List[Dict[str, Any]]]) -> Tuple[List[str], np.ndarray, np.ndarray]:
        payments = raw_data.get("payments", [])

        feature_names = [
            "amount", "time_to_pay_sec", "is_failed_status", "is_manual_transfer"
        ]

        feature_rows = []
        target_rows = []

        for pay in payments:
            amt = float(pay.get("amount", 0.0))
            status = pay.get("status", "paid")
            method = pay.get("method_id", "")

            created_dt = datetime.strptime(pay["created_at"], "%Y-%m-%d %H:%M:%S")
            paid_str = pay.get("paid_at")
            if paid_str:
                paid_dt = datetime.strptime(paid_str, "%Y-%m-%d %H:%M:%S")
                time_to_pay = (paid_dt - created_dt).total_seconds()
            else:
                time_to_pay = 86400.0 # 24 hour penalty

            is_failed = 1.0 if status == "failed" else 0.0
            is_manual = 1.0 if method == "MANUAL_TRANSFER" else 0.0

            row = [amt, time_to_pay, is_failed, is_manual]

            # Target binary anomaly flag
            is_anomaly = 1.0 if (is_failed == 1.0 or time_to_pay > 36000.0) else 0.0

            feature_rows.append(row)
            target_rows.append(is_anomaly)

        X = np.array(feature_rows, dtype=np.float32) if feature_rows else np.empty((0, len(feature_names)), dtype=np.float32)
        y = np.array(target_rows, dtype=np.float32) if target_rows else np.empty((0,), dtype=np.float32)

        logger.info(f"Generated Anomaly Detection dataset: X.shape={X.shape}, y.shape={y.shape}")
        return feature_names, X, y


class CourierSLAFeatureBuilder:
    """
    Builds Courier SLA & Delivery Duration feature dataset (Phase 2.4).
    Target: delivery_duration_hours (continuous float).
    Features: courier one-hot encodings, shipping_cost, dispatch_day_of_week, dispatch_hour, dispatch_is_weekend.
    """
    def build_features(self, raw_data: Dict[str, Any]) -> Tuple[List[str], np.ndarray, np.ndarray]:
        logger.info("Building Courier SLA & Delivery Duration features...")
        shipments = raw_data.get("shipments", [])

        feature_names = [
            "courier_jne",
            "courier_jnt",
            "courier_sicepat",
            "shipping_cost",
            "dispatch_day_of_week",
            "dispatch_hour",
            "dispatch_is_weekend"
        ]

        feature_rows = []
        target_rows = []

        for shp in shipments:
            courier = shp.get("courier_name", "JNE").upper()
            cost = float(shp.get("shipping_cost", 15000.0))
            created_str = shp.get("created_at")
            delivered_str = shp.get("delivered_at")

            if not created_str or not delivered_str:
                continue

            created_dt = datetime.strptime(created_str, "%Y-%m-%d %H:%M:%S")
            delivered_dt = datetime.strptime(delivered_str, "%Y-%m-%d %H:%M:%S")

            duration_hours = (delivered_dt - created_dt).total_seconds() / 3600.0
            if duration_hours < 0:
                continue

            is_jne = 1.0 if courier == "JNE" else 0.0
            is_jnt = 1.0 if courier == "JNT" else 0.0
            is_sicepat = 1.0 if courier == "SICEPAT" else 0.0
            dow = float(created_dt.weekday())
            hour = float(created_dt.hour)
            is_weekend = 1.0 if dow >= 5 else 0.0

            row = [is_jne, is_jnt, is_sicepat, cost, dow, hour, is_weekend]
            feature_rows.append(row)
            target_rows.append(duration_hours)

        X = np.array(feature_rows, dtype=np.float32) if feature_rows else np.empty((0, len(feature_names)), dtype=np.float32)
        y = np.array(target_rows, dtype=np.float32) if target_rows else np.empty((0,), dtype=np.float32)

        logger.info(f"Generated Courier SLA dataset: X.shape={X.shape}, y.shape={y.shape}")
        return feature_names, X, y


class StandardScaler:
    """Standardize features by removing the mean and scaling to unit variance."""
    def __init__(self):
        self.mean = None
        self.std = None

    def fit_transform(self, X: np.ndarray) -> np.ndarray:
        if X.size == 0:
            return X
        self.mean = np.mean(X, axis=0, keepdims=True)
        self.std = np.std(X, axis=0, keepdims=True) + 1e-8
        return (X - self.mean) / self.std

    def transform(self, X: np.ndarray) -> np.ndarray:
        if self.mean is None or self.std is None:
            raise ValueError("StandardScaler must be fitted before calling transform.")
        return (X - self.mean) / self.std

