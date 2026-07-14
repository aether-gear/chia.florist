import React from 'react';
import ditheredGarden from '@/assets/dithered-garden.png';

interface AuthLayoutProps {
  children: React.ReactNode;
}

export default function AuthLayout({ children }: AuthLayoutProps) {
  return (
    <div className="min-h-screen grid grid-cols-1 lg:grid-cols-12 bg-white font-sans antialiased text-zinc-900 selection:bg-emerald-100 selection:text-emerald-900">
      {/* Left panel: Auth Form and Info */}
      <div className="
        flex flex-col justify-between p-6 min-h-screen
        sm:p-10
        md:p-12
        lg:col-span-7
        xl:col-span-6
      ">
        {/* Brand Header */}
        <div className="flex items-center gap-2">
          {/* Creative, minimalist seed/sprout star brand icon */}
          <div className="w-9 h-9 rounded-xl bg-emerald-50 flex items-center justify-center border border-emerald-100/50">
            <svg className="w-5 h-5 text-emerald-600" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
              <path d="M12 22c0-5.523-4.477-10-10-10 5.523 0 10-4.477 10-10 0 5.523 4.477 10 10 10-5.523 0-10 4.477-10 10z" />
            </svg>
          </div>
          <span className="font-display font-semibold text-base tracking-tight text-emerald-950">
            Chia Florist
          </span>
        </div>

        {/* Dynamic Inner Form Container */}
        <div className="my-auto py-8 w-full max-w-[360px] mx-auto animate-in fade-in slide-in-from-bottom-4 duration-500 ease-out">
          {children}
        </div>

        {/* Footer */}
        <div className="flex items-center gap-4 text-xs font-medium text-zinc-400">
          <a href="#" className="hover:text-emerald-700 transition-colors">Help</a>
          <span className="text-zinc-200">/</span>
          <a href="#" className="hover:text-emerald-700 transition-colors">Terms</a>
          <span className="text-zinc-200">/</span>
          <a href="#" className="hover:text-emerald-700 transition-colors">Privacy</a>
        </div>
      </div>

      {/* Right panel: Full-screen 1-bit dithered art */}
      <div className="
        hidden bg-zinc-950 relative overflow-hidden select-none border-l border-zinc-100
        lg:block lg:col-span-5
        xl:col-span-6
      ">
        {/* High-contrast retro image */}
        <img
          src={ditheredGarden}
          alt="Florist Greenhouse Garden"
          className="absolute inset-0 w-full h-full object-cover filter brightness-[1.02] contrast-125"
          style={{ imageRendering: 'pixelated' }}
        />

        {/* Subtle, premium organic vignette */}
        <div className="absolute inset-0 bg-gradient-to-tr from-emerald-950/20 via-transparent to-transparent pointer-events-none" />
      </div>
    </div>
  );
}
