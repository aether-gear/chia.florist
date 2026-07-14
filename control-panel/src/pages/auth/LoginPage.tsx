import { Link, useNavigate } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Checkbox } from '@/components/ui/checkbox';
import { Loader2 } from 'lucide-react';
import { useState } from 'react';
import { fetchApi } from '@/lib/api';
import AuthLayout from '@/layouts/AuthLayout';

export default function LoginPage() {
  const navigate = useNavigate();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [rememberMe, setRememberMe] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setLoading(true);

    try {
      await fetchApi('/auth/staff/signin', {
        method: 'POST',
        body: JSON.stringify({
          email,
          password,
          user_agent: navigator.userAgent,
          ip_address: '127.0.0.1', // In a real app, backend determines this or use a 3rd party service
          remember_me: rememberMe
        })
      });

      // Clear both first to avoid mixed states
      localStorage.removeItem('isAuthenticated');
      localStorage.removeItem('userEmail');
      sessionStorage.removeItem('isAuthenticated');
      sessionStorage.removeItem('userEmail');

      const storage = rememberMe ? localStorage : sessionStorage;
      storage.setItem('isAuthenticated', 'true');
      storage.setItem('userEmail', email);
      navigate('/');
    } catch (err: any) {
      setError(err.message || 'Login failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <AuthLayout>
      <div className="space-y-6">
        {/* Title and Subtitle */}
        <div className="space-y-2">
          <h1 className="font-display text-3xl font-semibold tracking-tight text-emerald-950">
            Sign In
          </h1>
          <p className="text-sm text-zinc-500 font-sans font-light">
            Enter your credentials to access the
            <span className="text-emerald-950 italic"> Control Panel.</span>
          </p>
        </div>

        {/* Login Form */}
        <form onSubmit={handleLogin} className="space-y-4">
          {/* Email input field */}
          <div className="space-y-2">
            <label
              htmlFor="email"
              className="text-xs font-semibold tracking-wide text-zinc-500 uppercase font-sans"
            >
              Email address
            </label>
            <Input
              id="email"
              type="email"
              placeholder="admin@chia.florist"
              required
              onChange={(e) => setEmail(e.target.value)}
              className="w-full h-11 border-zinc-200/80 rounded-lg text-sm bg-zinc-50/20 focus-visible:bg-white focus:outline-none focus:ring-2 transition-all duration-200 font-sans"
            />
          </div>

          {/* Password input field */}
          <div className="space-y-2">
            <div className="flex items-center justify-between">
              <label
                htmlFor="password"
                className="text-xs font-semibold tracking-wide text-zinc-500 uppercase font-sans"
              >
                Password
              </label>
              <Link
                to="/forgot-password"
                className="text-xs font-semibold text-emerald-700 hover:text-emerald-800 transition-colors font-sans"
              >
                Forgot password?
              </Link>
            </div>
            <Input
              id="password"
              type="password"
              placeholder="••••••••"
              required
              onChange={(e) => setPassword(e.target.value)}
              className="w-full h-11 border-zinc-200/80 rounded-lg text-sm bg-zinc-50/20 focus-visible:bg-white focus:outline-none focus:ring-2 transition-all duration-200 font-sans"
            />
          </div>

          {/* Remember me option */}
          <div className="flex items-center space-x-2.5 py-1">
            <Checkbox
              id="rememberMe"
              checked={rememberMe}
              onCheckedChange={(checked) => setRememberMe(checked === true)}
              className="w-4 h-4 rounded border-zinc-300 focus-visible:ring-emerald-600"
            />
            <label
              htmlFor="rememberMe"
              className="text-xs font-medium text-zinc-500 cursor-pointer select-none font-sans"
            >
              Keep me signed in
            </label>
          </div>

          {/* Error message */}
          {error && (
            <div className="p-3 text-xs font-medium text-rose-600 bg-rose-50 border border-rose-100 rounded-lg font-sans animate-in fade-in zoom-in-95 duration-200">
              {error}
            </div>
          )}

          {/* Submit button */}
          <Button
            type="submit"
            disabled={loading}
            className="w-full h-11 bg-emerald-700 hover:bg-emerald-800 active:scale-[0.98] text-white font-medium rounded-lg shadow-sm transition-all duration-200 flex items-center justify-center font-sans mt-6"
          >
            {loading ? (
              <Loader2 className="h-4 w-4 animate-spin text-white" />
            ) : (
              'Sign In'
            )}
          </Button>
        </form>
      </div>
    </AuthLayout>
  );
}
