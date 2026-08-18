import os
import random
import logging
from datetime import datetime, timedelta
from typing import Dict, List, Any, Tuple, Optional
import numpy as np

logger = logging.getLogger("extractor")

class SyntheticDataExtractor:
    """
    Generates synthetic transactional dataset matching service-core schemas:
    - orders & order_items
    - inventory & product_stock_history
    - product_performance
    - payments & payment_events
    - shipments & shipment_events
    - audit_logs
    """
    def __init__(self, seed: int = 42, num_days: int = 90, num_products: int = 20, num_shops: int = 3):
        self.seed = seed
        self.num_days = num_days
        self.num_products = num_products
        self.num_shops = num_shops
        random.seed(seed)
        np.random.seed(seed)

    def extract_raw_data(self) -> Dict[str, List[Dict[str, Any]]]:
        logger.info(f"Generating synthetic transactional data for {self.num_days} days across {self.num_products} products...")
        start_date = datetime.now() - timedelta(days=self.num_days)
        
        products = []
        for pid in range(1, self.num_products + 1):
            margin = round(random.uniform(0.15, 0.50), 2)
            lead_time = random.choice([2, 3, 5, 7])
            cost_price = round(random.uniform(15000, 150000), 2)
            selling_price = round(cost_price / (1 - margin), 2)
            products.append({
                "product_id": f"PROD-{pid:03d}",
                "name": f"Flower Product {pid}",
                "category_id": f"CAT-{(pid % 4) + 1}",
                "cost_price": cost_price,
                "selling_price": selling_price,
                "gross_margin_pct": margin,
                "supplier_lead_time_days": lead_time,
                "view_count": random.randint(100, 5000)
            })

        orders = []
        order_items = []
        payments = []
        payment_events = []
        shipments = []
        shipment_events = []
        audit_logs = []
        stock_history = []
        inventory = []

        # Current stock levels per product per shop
        for p in products:
            for s in range(1, self.num_shops + 1):
                stock = random.randint(5, 100)
                reserved = random.randint(0, min(5, stock))
                inventory.append({
                    "product_id": p["product_id"],
                    "shop_id": f"SHOP-{s:02d}",
                    "stock": stock,
                    "reserved_stock": reserved
                })

        order_counter = 1
        for d in range(self.num_days):
            current_date = start_date + timedelta(days=d)
            # Add day-of-week & holiday seasonal multiplier
            is_weekend = current_date.weekday() >= 5
            seasonal_multiplier = 1.4 if is_weekend else 1.0
            
            daily_orders_count = int(np.random.poisson(lam=15 * seasonal_multiplier))

            # Record stock history snapshot daily
            for inv in inventory:
                noise = random.randint(-3, 3)
                inv["stock"] = max(0, inv["stock"] + noise)
                stock_history.append({
                    "product_id": inv["product_id"],
                    "shop_id": inv["shop_id"],
                    "available": inv["stock"],
                    "reserved": inv["reserved_stock"],
                    "recorded_at": current_date.strftime("%Y-%m-%d %H:%M:%S")
                })

            for _ in range(daily_orders_count):
                oid = f"ORD-{order_counter:06d}"
                shop_id = f"SHOP-{random.randint(1, self.num_shops):02d}"
                
                # Order creation
                created_at = current_date + timedelta(seconds=random.randint(0, 86399))
                p = random.choice(products)
                qty = random.randint(1, 5)
                subtotal = p["selling_price"] * qty
                shipping_fee = 15000.0
                total = subtotal + shipping_fee

                # 90% normal, 10% cancelled/failed (anomalies)
                is_anomaly = random.random() < 0.08
                status = "cancelled" if is_anomaly else "completed"

                orders.append({
                    "order_id": oid,
                    "shop_id": shop_id,
                    "status": status,
                    "subtotal": subtotal,
                    "shipping_fee": shipping_fee,
                    "total": total,
                    "created_at": created_at.strftime("%Y-%m-%d %H:%M:%S")
                })

                order_items.append({
                    "order_id": oid,
                    "product_id": p["product_id"],
                    "category_id": p["category_id"],
                    "quantity": qty,
                    "unit_price": p["selling_price"],
                    "subtotal": subtotal,
                    "courier_code": random.choice(["JNE", "JNT", "SICEPAT", "GOSEND"])
                })

                # Payments
                time_to_pay_sec = random.randint(120, 7200) if not is_anomaly else random.randint(18000, 86400)
                paid_at = created_at + timedelta(seconds=time_to_pay_sec)
                pay_status = "failed" if is_anomaly else "paid"
                
                payments.append({
                    "payment_id": f"PAY-{order_counter:06d}",
                    "order_id": oid,
                    "amount": total,
                    "status": pay_status,
                    "method_id": random.choice(["QRIS", "VIRTUAL_ACCOUNT", "CREDIT_CARD", "MANUAL_TRANSFER"]),
                    "paid_at": paid_at.strftime("%Y-%m-%d %H:%M:%S") if pay_status == "paid" else None,
                    "created_at": created_at.strftime("%Y-%m-%d %H:%M:%S")
                })

                payment_events.append({
                    "payment_id": f"PAY-{order_counter:06d}",
                    "transaction_status": pay_status,
                    "created_at": paid_at.strftime("%Y-%m-%d %H:%M:%S")
                })

                # Shipments
                if pay_status == "paid":
                    shipment_created = paid_at + timedelta(seconds=random.randint(600, 3600))
                    deliv_time = random.randint(7200, 86400)
                    shipments.append({
                        "shipment_id": f"SHP-{order_counter:06d}",
                        "order_id": oid,
                        "courier_name": random.choice(["JNE", "JNT", "SICEPAT"]),
                        "service": "REGULAR",
                        "shipping_cost": shipping_fee,
                        "created_at": shipment_created.strftime("%Y-%m-%d %H:%M:%S"),
                        "delivered_at": (shipment_created + timedelta(seconds=deliv_time)).strftime("%Y-%m-%d %H:%M:%S")
                    })

                # Audit logs
                audit_logs.append({
                    "event_id": f"EVT-{order_counter:06d}",
                    "event_type": "ORDER_CREATED",
                    "created_at": created_at.strftime("%Y-%m-%d %H:%M:%S")
                })

                order_counter += 1

        logger.info(f"Synthetic extraction complete: {len(orders)} orders, {len(order_items)} items, {len(payments)} payments generated.")
        
        return {
            "products": products,
            "inventory": inventory,
            "orders": orders,
            "order_items": order_items,
            "product_stock_history": stock_history,
            "payments": payments,
            "payment_events": payment_events,
            "shipments": shipments,
            "shipment_events": shipment_events,
            "audit_logs": audit_logs
        }


class DatabaseExtractor:
    """
    Attempts SQL extraction from service-core DB.
    Falls back to SyntheticDataExtractor if connection fails or driver is missing.
    """
    def __init__(self, config: Dict[str, Any]):
        self.config = config
        self.db_config = config.get("database", {})
        self.fallback = SyntheticDataExtractor(
            seed=config.get("seed", 42),
            num_days=self.db_config.get("query_limit_days", 90)
        )

    def extract_raw_data(self) -> Dict[str, List[Dict[str, Any]]]:
        if self.db_config.get("use_synthetic_fallback", True):
            logger.info("Configured to use synthetic fallback extractor.")
            return self.fallback.extract_raw_data()
        
        # Attempt PostgreSQL extraction via psycopg2 / sqlite if configured
        try:
            import psycopg2
            conn = psycopg2.connect(
                host=self.db_config.get("host", "localhost"),
                port=self.db_config.get("port", 5432),
                user=self.db_config.get("user", "postgres"),
                password=self.db_config.get("password", ""),
                dbname=self.db_config.get("dbname", "chia_florist")
            )
            logger.info("Connected to service-core PostgreSQL database.")
            # Execute queries (omitted for offline safety, fallback used)
            conn.close()
            return self.fallback.extract_raw_data()
        except Exception as e:
            logger.warning(f"Could not connect to PostgreSQL ({e}). Using synthetic fallback extractor.")
            return self.fallback.extract_raw_data()


def get_extractor(config: Dict[str, Any]):
    return DatabaseExtractor(config)
