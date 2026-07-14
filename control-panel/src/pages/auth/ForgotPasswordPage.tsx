import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Loader2, ArrowLeft, CheckCircle2 } from 'lucide-react';
import { fetchApi } from '@/lib/api';
import AuthLayout from '@/layouts/AuthLayout';

type Step = 'request' | 'verify' | 'reset' | 'success';

export default function ForgotPasswordPage() {
  const navigate = useNavigate();
  const [step, setStep] = useState<Step>('request');
  const [email, setEmail] = useState('');
  const [challengeId, setChallengeId] = useState('');
  const [otp, setOtp] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleRequestReset = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setLoading(true);

    try {
      const res = await fetchApi('/auth/staff/forgot-password', {
        method: 'POST',
        body: JSON.stringify({ email })
      });
      if (res && res.challenge_id) {
        setChallengeId(res.challenge_id);
      }
      setStep('verify');
    } catch (err: any) {
      setError(err.message || 'Failed to request password reset code.');
    } finally {
      setLoading(false);
    }
  };

  const handleVerifyOTP = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setLoading(true);

    try {
      const res = await fetchApi('/auth/forgot-password/verify', {
        method: 'POST',
        body: JSON.stringify({
          challenge_id: challengeId,
          otp: otp.trim()
        })
      });
      if (res && res.challenge_id) {
        setChallengeId(res.challenge_id);
      }
      setStep('reset');
    } catch (err: any) {
      setError(err.message || 'Invalid or expired verification code.');
    } finally {
      setLoading(false);
    }
  };

  const handleResetPassword = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    if (newPassword !== confirmPassword) {
      setError('Passwords do not match.');
      return;
    }

    if (newPassword.length < 8) {
      setError('Password must be at least 8 characters long.');
      return;
    }

    setLoading(true);

    try {
      await fetchApi('/auth/forgot-password/reset', {
        method: 'POST',
        body: JSON.stringify({
          challenge_id: challengeId,
          new_password: newPassword
        })
      });
      setStep('success');
    } catch (err: any) {
      setError(err.message || 'Failed to reset password.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <AuthLayout>
      <div className="space-y-6">
        {/* Title and Subtitle depending on step */}
        <div className="space-y-2">
          <h1 className="font-display text-3xl font-semibold tracking-tight text-emerald-950">
            {step === 'request' && 'Reset password'}
            {step === 'verify' && 'Enter reset code'}
            {step === 'reset' && 'Set new password'}
            {step === 'success' && 'Password reset'}
          </h1>
          <p className="text-sm text-zinc-500 font-sans font-light">
            {step === 'request' && 'Enter your email and we will send you a reset code.'}
            {step === 'verify' && `We sent a verification code to ${email}.`}
            {step === 'reset' && 'Choose a strong new password for your account.'}
            {step === 'success' && 'Your password has been successfully updated.'}
          </p>
        </div>

        {/* Global Error Display */}
        {error && (
          <div className="p-3 text-xs font-medium text-rose-600 bg-rose-50 border border-rose-100 rounded-lg font-sans animate-in fade-in zoom-in-95 duration-200">
            {error}
          </div>
        )}

        {/* Step 1: Request Reset */}
        {step === 'request' && (
          <form onSubmit={handleRequestReset} className="space-y-4">
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
                placeholder="name@company.com"
                required
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className="w-full h-11 border-zinc-200/80 rounded-lg text-sm bg-zinc-50/20 focus-visible:bg-white focus:outline-none focus:ring-2 transition-all duration-200 font-sans"
              />
            </div>
            <Button 
              type="submit" 
              disabled={loading} 
              className="w-full h-11 bg-emerald-700 hover:bg-emerald-800 active:scale-[0.98] text-white font-medium rounded-lg shadow-sm transition-all duration-200 flex items-center justify-center font-sans mt-6"
            >
              {loading ? <Loader2 className="h-4 w-4 animate-spin text-white" /> : 'Send reset code'}
            </Button>
          </form>
        )}

        {/* Step 2: Verify Code */}
        {step === 'verify' && (
          <form onSubmit={handleVerifyOTP} className="space-y-4">
            <div className="space-y-2">
              <label 
                htmlFor="otp" 
                className="text-xs font-semibold tracking-wide text-zinc-500 uppercase font-sans"
              >
                Verification Code
              </label>
              <Input
                id="otp"
                type="text"
                placeholder="Enter 6-digit code"
                required
                maxLength={6}
                value={otp}
                onChange={(e) => setOtp(e.target.value.replace(/\D/g, ''))}
                className="w-full h-11 border-zinc-200/80 rounded-lg text-center tracking-widest text-lg font-bold bg-zinc-50/20 focus-visible:bg-white focus:outline-none focus:ring-2 transition-all duration-200 font-sans"
              />
            </div>
            <Button 
              type="submit" 
              disabled={loading} 
              className="w-full h-11 bg-emerald-700 hover:bg-emerald-800 active:scale-[0.98] text-white font-medium rounded-lg shadow-sm transition-all duration-200 flex items-center justify-center font-sans mt-6"
            >
              {loading ? <Loader2 className="h-4 w-4 animate-spin text-white" /> : 'Verify code'}
            </Button>
          </form>
        )}

        {/* Step 3: Set New Password */}
        {step === 'reset' && (
          <form onSubmit={handleResetPassword} className="space-y-4">
            <div className="space-y-2">
              <label 
                htmlFor="newPassword" 
                className="text-xs font-semibold tracking-wide text-zinc-500 uppercase font-sans"
              >
                New Password
              </label>
              <Input
                id="newPassword"
                type="password"
                placeholder="••••••••"
                required
                value={newPassword}
                onChange={(e) => setNewPassword(e.target.value)}
                className="w-full h-11 border-zinc-200/80 rounded-lg text-sm bg-zinc-50/20 focus-visible:bg-white focus:outline-none focus:ring-2 transition-all duration-200 font-sans"
              />
            </div>
            <div className="space-y-2">
              <label 
                htmlFor="confirmPassword" 
                className="text-xs font-semibold tracking-wide text-zinc-500 uppercase font-sans"
              >
                Confirm Password
              </label>
              <Input
                id="confirmPassword"
                type="password"
                placeholder="••••••••"
                required
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                className="w-full h-11 border-zinc-200/80 rounded-lg text-sm bg-zinc-50/20 focus-visible:bg-white focus:outline-none focus:ring-2 transition-all duration-200 font-sans"
              />
            </div>
            <Button 
              type="submit" 
              disabled={loading} 
              className="w-full h-11 bg-emerald-700 hover:bg-emerald-800 active:scale-[0.98] text-white font-medium rounded-lg shadow-sm transition-all duration-200 flex items-center justify-center font-sans mt-6"
            >
              {loading ? <Loader2 className="h-4 w-4 animate-spin text-white" /> : 'Reset password'}
            </Button>
          </form>
        )}

        {/* Step 4: Success Message */}
        {step === 'success' && (
          <div className="space-y-6">
            <div className="flex items-center gap-3 text-emerald-700 bg-emerald-50 border border-emerald-100 rounded-xl p-4 animate-in fade-in zoom-in-95 duration-300">
              <CheckCircle2 className="h-5 w-5 shrink-0" />
              <span className="text-sm font-medium">Your password has been reset. You can now log in.</span>
            </div>
            <Button 
              onClick={() => navigate('/login')} 
              className="w-full h-11 bg-emerald-700 hover:bg-emerald-800 active:scale-[0.98] text-white font-medium rounded-lg shadow-sm transition-all duration-200 flex items-center justify-center font-sans"
            >
              Back to Sign In
            </Button>
          </div>
        )}

        {/* Bottom Back to Login Link */}
        {step !== 'success' && (
          <div className="pt-2">
            <Link 
              to="/login" 
              className="inline-flex items-center gap-1.5 text-xs font-semibold text-emerald-700 hover:text-emerald-800 transition-colors font-sans"
            >
              <ArrowLeft className="h-3.5 w-3.5" /> Back to Sign In
            </Link>
          </div>
        )}
      </div>
    </AuthLayout>
  );
}
