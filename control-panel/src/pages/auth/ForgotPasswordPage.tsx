import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardHeader, CardContent, CardFooter, CardTitle, CardDescription } from '@/components/ui/card';
import { KeyRound, Loader2, CheckCircle2, ArrowLeft, Mail, ShieldCheck } from 'lucide-react';
import { fetchApi } from '@/lib/api';

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
    <div className="min-h-screen flex items-center justify-center bg-slate-50 p-4">
      <Card className="w-full max-w-md shadow-lg border-0">
        <CardHeader className="space-y-3 items-center text-center pb-8">
          <div className="h-12 w-12 rounded-full bg-indigo-100 flex items-center justify-center mb-2">
            {step === 'success' ? (
              <CheckCircle2 className="h-7 w-7 text-emerald-600" />
            ) : (
              <KeyRound className="h-7 w-7 text-indigo-600" />
            )}
          </div>
          <CardTitle className="text-2xl font-bold">
            {step === 'request' && 'Reset password'}
            {step === 'verify' && 'Enter reset code'}
            {step === 'reset' && 'Set new password'}
            {step === 'success' && 'Password reset'}
          </CardTitle>
          <CardDescription>
            {step === 'request' && 'Enter your email and we will send you a reset code'}
            {step === 'verify' && `We sent a verification code to ${email}`}
            {step === 'reset' && 'Choose a strong new password for your account'}
            {step === 'success' && 'Your password has been successfully updated'}
          </CardDescription>
        </CardHeader>
        <CardContent>
          {error && (
            <div className="p-3 text-sm text-red-500 bg-red-50 rounded-md mb-4 border border-red-100">
              {error}
            </div>
          )}

          {step === 'request' && (
            <form onSubmit={handleRequestReset} className="space-y-4">
              <div className="space-y-2">
                <label className="text-sm font-medium leading-none" htmlFor="email">
                  Email address
                </label>
                <div className="relative">
                  <Input
                    id="email"
                    type="email"
                    placeholder="name@company.com"
                    required
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    className="pl-10"
                  />
                  <Mail className="absolute left-3 top-2.5 h-5 w-5 text-slate-400" />
                </div>
              </div>
              <Button type="submit" disabled={loading} className="w-full bg-indigo-600 hover:bg-indigo-700">
                {loading ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : 'Send reset code'}
              </Button>
            </form>
          )}

          {step === 'verify' && (
            <form onSubmit={handleVerifyOTP} className="space-y-4">
              <div className="space-y-2">
                <label className="text-sm font-medium leading-none" htmlFor="otp">
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
                  className="text-center tracking-widest text-lg font-bold"
                />
              </div>
              <Button type="submit" disabled={loading} className="w-full bg-indigo-600 hover:bg-indigo-700">
                {loading ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : 'Verify code'}
              </Button>
            </form>
          )}

          {step === 'reset' && (
            <form onSubmit={handleResetPassword} className="space-y-4">
              <div className="space-y-2">
                <label className="text-sm font-medium leading-none" htmlFor="newPassword">
                  New Password
                </label>
                <Input
                  id="newPassword"
                  type="password"
                  placeholder="••••••••"
                  required
                  value={newPassword}
                  onChange={(e) => setNewPassword(e.target.value)}
                />
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium leading-none" htmlFor="confirmPassword">
                  Confirm Password
                </label>
                <Input
                  id="confirmPassword"
                  type="password"
                  placeholder="••••••••"
                  required
                  value={confirmPassword}
                  onChange={(e) => setConfirmPassword(e.target.value)}
                />
              </div>
              <Button type="submit" disabled={loading} className="w-full bg-indigo-600 hover:bg-indigo-700">
                {loading ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : 'Reset password'}
              </Button>
            </form>
          )}

          {step === 'success' && (
            <div className="flex flex-col items-center justify-center py-4 space-y-4">
              <div className="flex items-center gap-2 text-emerald-600 bg-emerald-50 border border-emerald-100 rounded-lg p-3 w-full justify-center">
                <ShieldCheck className="h-5 w-5" />
                <span className="text-sm font-medium">Ready to sign in</span>
              </div>
              <Button onClick={() => navigate('/login')} className="w-full bg-indigo-600 hover:bg-indigo-700">
                Back to Sign In
              </Button>
            </div>
          )}
        </CardContent>
        {step !== 'success' && (
          <CardFooter className="flex justify-center border-t p-4 mt-2">
            <Link to="/login" className="flex items-center gap-1 text-sm text-indigo-600 hover:text-indigo-500 font-medium">
              <ArrowLeft className="h-4 w-4" /> Back to Sign In
            </Link>
          </CardFooter>
        )}
      </Card>
    </div>
  );
}
