import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import DashboardLayout from '../layouts/DashboardLayout';
import LoginPage from '../pages/auth/LoginPage';
import DashboardPage from '../pages/dashboard/DashboardPage';
import SecurityPage from '../pages/security/SecurityPage';
import PlaceholderPage from '../pages/PlaceholderPage';
import MerchantProfileSettings from '../pages/merchant-profile/MerchantProfileSettings';
import CreateMerchantPage from '../pages/admin/CreateMerchantPage';
import AddMerchantAccountPage from '../pages/admin/AddMerchantAccountPage';

export default function AppRoutes() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<LoginPage />} />
        
        <Route element={<DashboardLayout />}>
          <Route path="/" element={<DashboardPage />} />
          <Route path="/security" element={<SecurityPage />} />
          <Route path="/shop" element={<PlaceholderPage title="Shop" />} />
          <Route path="/products" element={<PlaceholderPage title="Products" />} />
          <Route path="/orders" element={<PlaceholderPage title="Orders" />} />
          <Route path="/transactions" element={<PlaceholderPage title="Transactions" />} />
          <Route path="/shipments" element={<PlaceholderPage title="Shipments" />} />
          <Route path="/merchant/settings" element={<MerchantProfileSettings />} />
          
          {/* Admin Routes */}
          <Route path="/admin/merchants/create" element={<CreateMerchantPage />} />
          <Route path="/admin/merchants/accounts/add" element={<AddMerchantAccountPage />} />
          <Route path="/admin/merchants/:merchantId/accounts/add" element={<AddMerchantAccountPage />} />
        </Route>
        
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </BrowserRouter>
  );
}
