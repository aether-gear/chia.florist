<!-- app/components/CToastContainer.vue -->
<!-- Global toast renderer. Place <CToastContainer /> once in app.vue or a layout. -->
<!-- Mirrors the custom.vue success-toast design with the same animation. -->
<script setup lang="ts">
import { useToast } from '~/composables/useToast'
const { toasts, dismiss } = useToast()
</script>

<template>
  <Teleport to="body">
    <div class="toast-stack" aria-live="polite">
      <TransitionGroup name="toast-slide" tag="div" class="toast-group">
        <div
          v-for="toast in toasts"
          :key="toast.id"
          class="c-toast"
          :class="`c-toast--${toast.type}`"
          role="status"
          @click="dismiss(toast.id)"
        >
          <!-- Icon -->
          <div class="ct-icon">
            <!-- Success -->
            <svg v-if="toast.type === 'success'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round">
              <path d="M5 13l4 4L19 7"/>
            </svg>
            <!-- Error -->
            <svg v-else-if="toast.type === 'error'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round">
              <line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
            <!-- Info -->
            <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round">
              <circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/>
            </svg>
          </div>
          <!-- Text -->
          <div class="ct-body">
            <p class="ct-title">{{ toast.title }}</p>
            <p v-if="toast.subtitle" class="ct-sub">{{ toast.subtitle }}</p>
          </div>
          <!-- Dismiss -->
          <button class="ct-close" @click.stop="dismiss(toast.id)" aria-label="Dismiss">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
              <line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
          </button>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<style scoped>
/* ── Stack container ───────────────────────────────────────────────── */
.toast-stack {
  position: fixed;
  bottom: 1.5rem;
  left: 50%;
  transform: translateX(-50%);
  z-index: 9999;
  display: flex;
  flex-direction: column;
  gap: 0.6rem;
  align-items: center;
  pointer-events: none;
}
.toast-group {
  display: flex;
  flex-direction: column;
  gap: 0.6rem;
  align-items: center;
}

/* ── Individual toast ─────────────────────────────────────────────── */
.c-toast {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  background: #fff;
  border: 1px solid #e5ddd4;
  border-radius: 12px;
  padding: 0.875rem 1rem 0.875rem 1.1rem;
  box-shadow: 0 8px 32px rgba(0,0,0,0.13);
  min-width: 260px;
  max-width: 360px;
  pointer-events: all;
  cursor: pointer;
  transition: box-shadow 0.2s;
}
.c-toast:hover { box-shadow: 0 12px 40px rgba(0,0,0,0.18); }

/* ── Type colour tokens ──────────────────────────────────────────── */
.ct-icon {
  width: 34px; height: 34px;
  border-radius: 50%;
  display: flex; align-items: center; justify-content: center;
  flex-shrink: 0;
}
.ct-icon svg { width: 1rem; height: 1rem; }

.c-toast--success .ct-icon {
  background: rgba(107,143,110,0.12);
  color: #6b8f6e;
  border: 1px solid rgba(107,143,110,0.3);
}
.c-toast--error .ct-icon {
  background: rgba(196,61,61,0.1);
  color: #c43d3d;
  border: 1px solid rgba(196,61,61,0.25);
}
.c-toast--info .ct-icon {
  background: rgba(59,130,246,0.1);
  color: #3b82f6;
  border: 1px solid rgba(59,130,246,0.25);
}

/* ── Text ─────────────────────────────────────────────────────────── */
.ct-body { flex: 1; min-width: 0; }
.ct-title { font-size: 0.82rem; font-weight: 800; color: #1c1813; margin-bottom: 0.1rem; }
.ct-sub   { font-size: 0.68rem; color: #a8998d; }

/* ── Close button ────────────────────────────────────────────────── */
.ct-close {
  flex-shrink: 0;
  width: 24px; height: 24px;
  border-radius: 50%;
  background: transparent;
  border: none;
  color: #b0a099;
  display: flex; align-items: center; justify-content: center;
  cursor: pointer;
  transition: background 0.15s, color 0.15s;
}
.ct-close:hover { background: #f4f0eb; color: #1c1813; }
.ct-close svg { width: 0.75rem; height: 0.75rem; }

/* ── Transition (same as custom.vue success-toast) ───────────────── */
.toast-slide-enter-active,
.toast-slide-leave-active {
  transition: opacity 0.3s ease, transform 0.35s cubic-bezier(0.34,1.56,0.64,1);
}
.toast-slide-enter-from,
.toast-slide-leave-to {
  opacity: 0;
  transform: translateY(18px) scale(0.95);
}

@media (max-width: 480px) {
  .c-toast { min-width: calc(100vw - 2rem); max-width: calc(100vw - 2rem); }
}
</style>
