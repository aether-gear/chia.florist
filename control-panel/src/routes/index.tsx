import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import DashboardLayout from '../layouts/DashboardLayout';
import LoginPage from '../pages/auth/LoginPage';
import ForgotPasswordPage from '../pages/auth/ForgotPasswordPage';
import DashboardPage from '../pages/dashboard/DashboardPage';
import SecurityPage from '../pages/security/SecurityPage';
import OrdersPage from '../pages/orders/OrdersPage';
import MerchantProfileSettings from '../pages/merchant-profile/MerchantProfileSettings';
import CreateMerchantPage from '../pages/admin/CreateMerchantPage';
import AddMerchantAccountPage from '../pages/admin/AddMerchantAccountPage';
import ProductsPage from '../pages/dashboard/ProductsPage';

import MerchantsListPage from '../pages/admin/MerchantsListPage';
import CustomersListPage from '../pages/admin/CustomersListPage';
import ShopManagementPage from '../pages/shop/ShopManagementPage';
import PaymentSettingsPage from '../pages/admin/payments/PaymentSettingsPage';
import CreatePaymentAccountPage from '../pages/admin/payments/CreatePaymentAccountPage';
import ProtectedRoute from '../components/ProtectedRoute';
import AuditLogsPage from '../pages/admin/AuditLogsPage';
import AnalyticsPage from '../pages/admin/analytics/AnalyticsPage';

export default function AppRoutes() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<LoginPage />} />
        <Route path="/forgot-password" element={<ForgotPasswordPage />} />
        
        <Route element={<ProtectedRoute />}>
          <Route element={<DashboardLayout />}>
            <Route path="/" element={<DashboardPage />} />
            <Route path="/security" element={<SecurityPage />} />
            <Route path="/shop" element={<ShopManagementPage />} />
            <Route path="/products" element={<ProductsPage />} />

            <Route path="/orders" element={<OrdersPage />} />
            <Route path="/shipments" element={<Navigate to="/orders" replace />} />
            <Route path="/merchant/settings" element={<MerchantProfileSettings />} />
            
            {/* Admin Routes */}
            <Route path="/admin/analytics" element={<AnalyticsPage />} />
            <Route path="/admin/merchants" element={<MerchantsListPage />} />
            <Route path="/admin/customers" element={<CustomersListPage />} />
            <Route path="/admin/payments" element={<PaymentSettingsPage />} />
            <Route path="/admin/payments/accounts/create" element={<CreatePaymentAccountPage />} />
            <Route path="/admin/merchants/create" element={<CreateMerchantPage />} />
            <Route path="/admin/merchants/accounts/add" element={<AddMerchantAccountPage />} />
            <Route path="/admin/merchants/:merchantId/accounts/add" element={<AddMerchantAccountPage />} />
            <Route path="/admin/audit-logs" element={<AuditLogsPage />} />
          </Route>
        </Route>
        
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </BrowserRouter>
  );
}

