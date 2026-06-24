import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import DashboardLayout from '../layouts/DashboardLayout';
import LoginPage from '../pages/auth/LoginPage';
import DashboardPage from '../pages/dashboard/DashboardPage';
import SecurityPage from '../pages/security/SecurityPage';
import PlaceholderPage from '../pages/PlaceholderPage';
import OrdersPage from '../pages/orders/OrdersPage';
import MerchantProfileSettings from '../pages/merchant-profile/MerchantProfileSettings';
import CreateMerchantPage from '../pages/admin/CreateMerchantPage';
import AddMerchantAccountPage from '../pages/admin/AddMerchantAccountPage';
import ProductsPage from '../pages/dashboard/ProductsPage';
import CreateProductPage from '../pages/products/CreateProductPage';
import MerchantsListPage from '../pages/admin/MerchantsListPage';
import CustomersListPage from '../pages/admin/CustomersListPage';
import ShopManagementPage from '../pages/shop/ShopManagementPage';
import PaymentMethodsPage from '../pages/admin/payments/PaymentMethodsPage';
import PaymentAccountsPage from '../pages/admin/payments/PaymentAccountsPage';
import CreatePaymentAccountPage from '../pages/admin/payments/CreatePaymentAccountPage';
import ProtectedRoute from '../components/ProtectedRoute';

export default function AppRoutes() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<LoginPage />} />
        
        <Route element={<ProtectedRoute />}>
          <Route element={<DashboardLayout />}>
            <Route path="/" element={<DashboardPage />} />
            <Route path="/security" element={<SecurityPage />} />
            <Route path="/shop" element={<ShopManagementPage />} />
            <Route path="/products" element={<ProductsPage />} />
            <Route path="/products/create" element={<CreateProductPage />} />
            <Route path="/orders" element={<OrdersPage />} />
            <Route path="/transactions" element={<PlaceholderPage title="Transactions" />} />
            <Route path="/shipments" element={<PlaceholderPage title="Shipments" />} />
            <Route path="/merchant/settings" element={<MerchantProfileSettings />} />
            
            {/* Admin Routes */}
            <Route path="/admin/merchants" element={<MerchantsListPage />} />
            <Route path="/admin/customers" element={<CustomersListPage />} />
            <Route path="/admin/payments/methods" element={<PaymentMethodsPage />} />
            <Route path="/admin/payments/accounts" element={<PaymentAccountsPage />} />
            <Route path="/admin/payments/accounts/create" element={<CreatePaymentAccountPage />} />
            <Route path="/admin/merchants/create" element={<CreateMerchantPage />} />
            <Route path="/admin/merchants/accounts/add" element={<AddMerchantAccountPage />} />
            <Route path="/admin/merchants/:merchantId/accounts/add" element={<AddMerchantAccountPage />} />
          </Route>
        </Route>
        
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </BrowserRouter>
  );
}

