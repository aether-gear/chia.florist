import { useState, useRef, useEffect } from 'react';
import { Camera, Trash2, Loader2 } from 'lucide-react';
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { useAvatarUpload } from '@/hooks/useAvatarUpload';

interface AvatarUploadProps {
  userId: string;
  currentAvatarUrl?: string;
  displayName?: string;
  onUploadComplete: (url: string) => void;
  onRemoveComplete: () => void;
}

type DragMode = 'img' | 'move' | 'nw' | 'n' | 'ne' | 'e' | 'se' | 's' | 'sw' | 'w' | null;

const CONT_W = 440;
const CONT_H = 320;
const MIN_CROP = 40;

export default function AvatarUpload({
  userId,
  currentAvatarUrl,
  displayName,
  onUploadComplete,
  onRemoveComplete,
}: AvatarUploadProps) {
  const { uploading, deleting, upload, remove } = useAvatarUpload();
  const fileInputRef = useRef<HTMLInputElement>(null);
  const cropImgRef = useRef<HTMLImageElement>(null);

  // Cropper dialog states
  const [showCropper, setShowCropper] = useState(false);
  const [cropSrc, setCropSrc] = useState('');
  const [cropFilename, setCropFilename] = useState('avatar.jpg');
  const [cropMime, setCropMime] = useState('image/jpeg');

  // Image natural dimensions
  const [naturalW, setNaturalW] = useState(0);
  const [naturalH, setNaturalH] = useState(0);

  // Zoom & Image Pan offsets
  const [cropZoom, setCropZoom] = useState(1);
  const [baseScale, setBaseScale] = useState(1);
  const [imgX, setImgX] = useState(0);
  const [imgY, setImgY] = useState(0);

  const scale = baseScale * cropZoom;
  const maxImgDX = Math.max(0, (naturalW * scale - CONT_W) / 2);
  const maxImgDY = Math.max(0, (naturalH * scale - CONT_H) / 2);

  // Crop selection box coordinates (in container pixel space)
  const [cbX, setCbX] = useState(60);
  const [cbY, setCbY] = useState(40);
  const [cbW, setCbW] = useState(CONT_W - 120);
  const [cbH, setCbH] = useState(CONT_H - 80);

  // Dragging State
  const dragMode = useRef<DragMode>(null);
  const dragStartX = useRef(0);
  const dragStartY = useRef(0);
  const initCbX = useRef(0);
  const initCbY = useRef(0);
  const initCbW = useRef(0);
  const initCbH = useRef(0);
  const lastPanX = useRef(0);
  const lastPanY = useRef(0);

  // Image load handler
  const onImgLoad = (e: React.SyntheticEvent<HTMLImageElement>) => {
    const img = e.currentTarget;
    setNaturalW(img.naturalWidth);
    setNaturalH(img.naturalHeight);

    const bScale = Math.max(CONT_W / img.naturalWidth, CONT_H / img.naturalHeight);
    setBaseScale(bScale);
    setCropZoom(1);
    setImgX(0);
    setImgY(0);

    // Square crop box, 75% of shorter container side, centered
    const side = Math.round(Math.min(CONT_W, CONT_H) * 0.75);
    setCbW(side);
    setCbH(side);
    setCbX(Math.round((CONT_W - side) / 2));
    setCbY(Math.round((CONT_H - side) / 2));
  };

  const startCropHandle = (e: React.MouseEvent | React.TouchEvent, mode: DragMode) => {
    e.stopPropagation();
    const clientX = 'touches' in e ? e.touches[0].clientX : e.clientX;
    const clientY = 'touches' in e ? e.touches[0].clientY : e.clientY;

    dragMode.current = mode;
    dragStartX.current = clientX;
    dragStartY.current = clientY;
    initCbX.current = cbX;
    initCbY.current = cbY;
    initCbW.current = cbW;
    initCbH.current = cbH;
  };

  const startImgPan = (e: React.MouseEvent | React.TouchEvent) => {
    const clientX = 'touches' in e ? e.touches[0].clientX : e.clientX;
    const clientY = 'touches' in e ? e.touches[0].clientY : e.clientY;

    dragMode.current = 'img';
    lastPanX.current = clientX;
    lastPanY.current = clientY;
  };

  const handleGlobalMove = (clientX: number, clientY: number) => {
    const m = dragMode.current;
    if (!m) return;

    if (m === 'img') {
      const dx = clientX - lastPanX.current;
      const dy = clientY - lastPanY.current;
      lastPanX.current = clientX;
      lastPanY.current = clientY;
      
      // We need to clamp based on the state update, so calculate next coordinates
      setImgX((prevX) => {
        const nextX = prevX + dx;
        return Math.max(-maxImgDX, Math.min(maxImgDX, nextX));
      });
      setImgY((prevY) => {
        const nextY = prevY + dy;
        return Math.max(-maxImgDY, Math.min(maxImgDY, nextY));
      });
      return;
    }

    const dx = clientX - dragStartX.current;
    const dy = clientY - dragStartY.current;
    let nx = initCbX.current;
    let ny = initCbY.current;
    let nw = initCbW.current;
    let nh = initCbH.current;

    if (m === 'move') {
      nx = Math.max(0, Math.min(CONT_W - nw, nx + dx));
      ny = Math.max(0, Math.min(CONT_H - nh, ny + dy));
    } else {
      let rawSide = nw;
      if (m === 'e') rawSide = initCbW.current + dx;
      if (m === 'w') rawSide = initCbW.current - dx;
      if (m === 's') rawSide = initCbH.current + dy;
      if (m === 'n') rawSide = initCbH.current - dy;
      if (m === 'se') rawSide = Math.max(initCbW.current + dx, initCbH.current + dy);
      if (m === 'sw') rawSide = Math.max(initCbW.current - dx, initCbH.current + dy);
      if (m === 'ne') rawSide = Math.max(initCbW.current + dx, initCbH.current - dy);
      if (m === 'nw') rawSide = Math.max(initCbW.current - dx, initCbH.current - dy);

      rawSide = Math.max(MIN_CROP, rawSide);
      nw = rawSide;
      nh = rawSide;

      if (m.includes('w')) nx = initCbX.current + initCbW.current - rawSide;
      if (m.includes('n')) ny = initCbY.current + initCbH.current - rawSide;

      nx = Math.max(0, nx);
      ny = Math.max(0, ny);
      if (nx + nw > CONT_W) {
        nw = CONT_W - nx;
        nh = nw;
      }
      if (ny + nh > CONT_H) {
        nh = CONT_H - ny;
        nw = nh;
      }
    }

    setCbX(nx);
    setCbY(ny);
    setCbW(nw);
    setCbH(nh);
  };

  const handleGlobalEnd = () => {
    dragMode.current = null;
  };

  // Add global move/end event listeners during crop session
  useEffect(() => {
    if (!showCropper) return;

    const onMouseMove = (e: MouseEvent) => {
      handleGlobalMove(e.clientX, e.clientY);
    };
    const onMouseUp = () => {
      handleGlobalEnd();
    };

    const onTouchMove = (e: TouchEvent) => {
      if (e.touches.length === 1) {
        // Prevent window bouncing/scroll during crop drag
        e.preventDefault();
        handleGlobalMove(e.touches[0].clientX, e.touches[0].clientY);
      }
    };
    const onTouchEnd = () => {
      handleGlobalEnd();
    };

    window.addEventListener('mousemove', onMouseMove);
    window.addEventListener('mouseup', onMouseUp);
    window.addEventListener('touchmove', onTouchMove, { passive: false });
    window.addEventListener('touchend', onTouchEnd);

    return () => {
      window.removeEventListener('mousemove', onMouseMove);
      window.removeEventListener('mouseup', onMouseUp);
      window.removeEventListener('touchmove', onTouchMove);
      window.removeEventListener('touchend', onTouchEnd);
    };
  }, [showCropper, maxImgDX, maxImgDY, cbX, cbY, cbW, cbH]);

  const onCropZoomChange = (newZoom: number) => {
    setCropZoom(newZoom);
    // Recalculate clamps using next state
    const nextScale = baseScale * newZoom;
    const nextMaxDX = Math.max(0, (naturalW * nextScale - CONT_W) / 2);
    const nextMaxDY = Math.max(0, (naturalH * nextScale - CONT_H) / 2);
    setImgX((prev) => Math.max(-nextMaxDX, Math.min(nextMaxDX, prev)));
    setImgY((prev) => Math.max(-nextMaxDY, Math.min(nextMaxDY, prev)));

    // Keep crop box within container
    setCbW((w) => {
      let nextW = w;
      setCbX((x) => {
        if (x + nextW > CONT_W) {
          nextW = CONT_W - x;
        }
        return x;
      });
      return nextW;
    });

    setCbH((h) => {
      let nextH = h;
      setCbY((y) => {
        if (y + nextH > CONT_H) {
          nextH = CONT_H - y;
        }
        return y;
      });
      return nextH;
    });
  };

  const cancelCrop = () => {
    setShowCropper(false);
    setCropSrc('');
    if (fileInputRef.current) fileInputRef.current.value = '';
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    if (!file.type.startsWith('image/')) {
      alert('Please select an image file.');
      return;
    }
    if (file.size > 10 * 1024 * 1024) {
      alert('Image must be less than 10 MB.');
      return;
    }

    setCropFilename(file.name || 'avatar.jpg');
    setCropMime(file.type || 'image/jpeg');

    const reader = new FileReader();
    reader.onload = (ev) => {
      if (ev.target?.result) {
        setCropSrc(ev.target.result as string);
        setShowCropper(true);
      }
    };
    reader.readAsDataURL(file);
  };

  const performCropAndUpload = async () => {
    const img = cropImgRef.current;
    if (!img) return;

    const canvas = document.createElement('canvas');
    canvas.width = 800;
    canvas.height = 800;
    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    const s = baseScale * cropZoom;
    const imgLeft = CONT_W / 2 - (naturalW * s) / 2 + imgX;
    const imgTop = CONT_H / 2 - (naturalH * s) / 2 + imgY;

    // Crop box in source image coordinates
    const srcX = (cbX - imgLeft) / s;
    const srcY = (cbY - imgTop) / s;
    const srcW = cbW / s;
    const srcH = cbH / s;

    ctx.drawImage(img, srcX, srcY, srcW, srcH, 0, 0, 800, 800);

    canvas.toBlob(async (blob) => {
      if (!blob) {
        alert('Crop failed.');
        return;
      }

      // Close modal first
      setShowCropper(false);
      setCropSrc('');
      if (fileInputRef.current) fileInputRef.current.value = '';

      const croppedFile = new File([blob], cropFilename, { type: cropMime });
      const urls = await upload(userId, croppedFile);
      if (urls) {
        onUploadComplete(urls.signedUrl || urls.publicUrl);
      }
    }, cropMime);
  };

  const handleRemovePicture = async () => {
    if (window.confirm('Are you sure you want to remove your profile picture?')) {
      const success = await remove(userId);
      if (success) {
        onRemoveComplete();
      }
    }
  };

  return (
    <div className="flex items-center space-x-6">
      {/* Hidden File Input */}
      <input
        type="file"
        ref={fileInputRef}
        onChange={handleFileChange}
        accept="image/*"
        className="hidden"
      />

      {/* Avatar Container */}
      <div className="relative group shrink-0">
        <div className="h-20 w-20 rounded-full bg-primary/10 overflow-hidden border border-border flex items-center justify-center relative">
          {currentAvatarUrl ? (
            <img src={currentAvatarUrl} alt="Profile" className="h-full w-full object-cover" />
          ) : (
            <span className="text-primary font-bold text-xl uppercase">
              {displayName?.charAt(0) || 'S'}
            </span>
          )}

          {/* Loading overlay when uploading or deleting */}
          {(uploading || deleting) && (
            <div className="absolute inset-0 bg-background/50 flex items-center justify-center">
              <Loader2 className="h-6 w-6 animate-spin text-primary" />
            </div>
          )}
        </div>

        {/* Change Photo Badge Button */}
        {!uploading && !deleting && (
          <button
            type="button"
            onClick={() => fileInputRef.current?.click()}
            title="Change photo"
            aria-label="Change photo"
            className="absolute bottom-0 right-0 h-7 w-7 rounded-full bg-primary text-primary-foreground flex items-center justify-center shadow-md hover:bg-primary/95 transition focus:outline-none ring-2 ring-background"
          >
            <Camera className="h-3.5 w-3.5" />
          </button>
        )}
      </div>

      {/* Remove Picture Action */}
      {currentAvatarUrl && !uploading && !deleting && (
        <Button
          type="button"
          variant="outline"
          onClick={handleRemovePicture}
          className="text-rose-600 border-rose-200 hover:bg-rose-50 rounded-xl flex items-center space-x-2 text-xs h-9 px-3"
        >
          <Trash2 className="h-3.5 w-3.5" />
          <span>Remove Photo</span>
        </Button>
      )}

      {/* Reusable Canvas Crop Modal */}
      {showCropper && (
        <Dialog open={showCropper} onOpenChange={(open) => !open && cancelCrop()}>
          <DialogContent className="max-w-[480px] p-0 overflow-hidden gap-0 rounded-3xl border-none">
            <DialogHeader className="px-6 pt-5 pb-3 border-b border-border">
              <DialogTitle className="font-bold font-display text-lg text-foreground">
                Crop Profile Picture
              </DialogTitle>
              <p className="text-[11px] text-muted-foreground mt-1 leading-relaxed">
                Drag corners to resize selection · Drag inside selection to move it · Drag outside to pan image
              </p>
            </DialogHeader>

            {/* Canvas Cropper Area */}
            <div
              className="relative bg-black overflow-hidden select-none"
              style={{ width: `${CONT_W}px`, height: `${CONT_H}px`, margin: '0 auto' }}
              onMouseDown={startImgPan}
              onTouchStart={startImgPan}
            >
              {/* Image */}
              {cropSrc && (
                <img
                  ref={cropImgRef}
                  src={cropSrc}
                  alt=""
                  onLoad={onImgLoad}
                  className="max-w-none max-h-none absolute pointer-events-none select-none"
                  style={{
                    width: `${naturalW * scale}px`,
                    height: `${naturalH * scale}px`,
                    left: '50%',
                    top: '50%',
                    transform: `translate(-50%, -50%) translate(${imgX}px, ${imgY}px)`,
                  }}
                  draggable={false}
                />
              )}

              {/* Dim Out of Crop Overlays */}
              <div
                className="absolute pointer-events-none bg-black/60"
                style={{ top: 0, left: 0, width: '100%', height: `${cbY}px` }}
              />
              <div
                className="absolute pointer-events-none bg-black/60"
                style={{ top: `${cbY + cbH}px`, left: 0, width: '100%', bottom: 0 }}
              />
              <div
                className="absolute pointer-events-none bg-black/60"
                style={{ top: `${cbY}px`, left: 0, width: `${cbX}px`, height: `${cbH}px` }}
              />
              <div
                className="absolute pointer-events-none bg-black/60"
                style={{ top: `${cbY}px`, left: `${cbX + cbW}px`, right: 0, height: `${cbH}px` }}
              />

              {/* Crop Selection Frame */}
              <div
                className="absolute border-2 border-white cursor-move box-border"
                style={{
                  left: `${cbX}px`,
                  top: `${cbY}px`,
                  width: `${cbW}px`,
                  height: `${cbH}px`,
                }}
                onMouseDown={(e) => startCropHandle(e, 'move')}
                onTouchStart={(e) => startCropHandle(e, 'move')}
              >
                {/* Rule of thirds grid lines */}
                <div
                  className="absolute inset-0 pointer-events-none"
                  style={{
                    backgroundImage:
                      'linear-gradient(rgba(255,255,255,0.2) 1px, transparent 1px), linear-gradient(90deg, rgba(255,255,255,0.2) 1px, transparent 1px)',
                    backgroundSize: '33.33% 33.33%',
                  }}
                />

                {/* Handles */}
                {/* Corners */}
                <div
                  className="absolute w-2.5 h-2.5 bg-white border border-primary rounded-[2px]"
                  style={{ top: -5, left: -5, cursor: 'nw-resize' }}
                  onMouseDown={(e) => startCropHandle(e, 'nw')}
                  onTouchStart={(e) => startCropHandle(e, 'nw')}
                />
                <div
                  className="absolute w-2.5 h-2.5 bg-white border border-primary rounded-[2px]"
                  style={{ top: -5, right: -5, cursor: 'ne-resize' }}
                  onMouseDown={(e) => startCropHandle(e, 'ne')}
                  onTouchStart={(e) => startCropHandle(e, 'ne')}
                />
                <div
                  className="absolute w-2.5 h-2.5 bg-white border border-primary rounded-[2px]"
                  style={{ bottom: -5, left: -5, cursor: 'sw-resize' }}
                  onMouseDown={(e) => startCropHandle(e, 'sw')}
                  onTouchStart={(e) => startCropHandle(e, 'sw')}
                />
                <div
                  className="absolute w-2.5 h-2.5 bg-white border border-primary rounded-[2px]"
                  style={{ bottom: -5, right: -5, cursor: 'se-resize' }}
                  onMouseDown={(e) => startCropHandle(e, 'se')}
                  onTouchStart={(e) => startCropHandle(e, 'se')}
                />

                {/* Edges */}
                <div
                  className="absolute w-6 h-1.5 bg-white border border-primary rounded-[1px]"
                  style={{ top: -4, left: 'calc(50% - 12px)', cursor: 'n-resize' }}
                  onMouseDown={(e) => startCropHandle(e, 'n')}
                  onTouchStart={(e) => startCropHandle(e, 'n')}
                />
                <div
                  className="absolute w-6 h-1.5 bg-white border border-primary rounded-[1px]"
                  style={{ bottom: -4, left: 'calc(50% - 12px)', cursor: 's-resize' }}
                  onMouseDown={(e) => startCropHandle(e, 's')}
                  onTouchStart={(e) => startCropHandle(e, 's')}
                />
                <div
                  className="absolute w-1.5 h-6 bg-white border border-primary rounded-[1px]"
                  style={{ left: -4, top: 'calc(50% - 12px)', cursor: 'w-resize' }}
                  onMouseDown={(e) => startCropHandle(e, 'w')}
                  onTouchStart={(e) => startCropHandle(e, 'w')}
                />
                <div
                  className="absolute w-1.5 h-6 bg-white border border-primary rounded-[1px]"
                  style={{ right: -4, top: 'calc(50% - 12px)', cursor: 'e-resize' }}
                  onMouseDown={(e) => startCropHandle(e, 'e')}
                  onTouchStart={(e) => startCropHandle(e, 'e')}
                />
              </div>
            </div>

            {/* Dialog Footer Controls */}
            <div className="px-6 py-4 space-y-4 border-t border-border">
              {/* Zoom Slider */}
              <div className="flex items-center gap-3">
                <span className="text-xs font-bold text-muted-foreground w-10">Zoom</span>
                <input
                  type="range"
                  min="1"
                  max="4"
                  step="0.01"
                  value={cropZoom}
                  onChange={(e) => onCropZoomChange(parseFloat(e.target.value))}
                  className="flex-1 h-1.5 rounded-full appearance-none cursor-pointer accent-primary bg-muted"
                />
                <span className="text-xs font-bold text-foreground w-10 text-right">
                  {Math.round(cropZoom * 100)}%
                </span>
              </div>

              {/* Bottom Actions */}
              <div className="flex items-center justify-between">
                <span className="text-[10px] text-muted-foreground">
                  ⬜ {cbW} × {cbH} px (square)
                </span>
                <div className="flex space-x-2">
                  <Button variant="outline" onClick={cancelCrop} className="rounded-xl text-xs h-9">
                    Cancel
                  </Button>
                  <Button onClick={performCropAndUpload} className="rounded-xl text-xs h-9 bg-primary hover:bg-primary/90 text-primary-foreground">
                    Crop & Save
                  </Button>
                </div>
              </div>
            </div>
          </DialogContent>
        </Dialog>
      )}
    </div>
  );
}
