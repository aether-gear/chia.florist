import { useState } from 'react';
import { supabaseStorage } from '@/lib/supabaseStorage';
import { useToast } from './use-toast';

export function useAvatarUpload() {
  const [uploading, setUploading] = useState(false);
  const [deleting, setDeleting] = useState(false);
  const [uploadError, setUploadError] = useState<string | null>(null);
  const { toast } = useToast();

  const upload = async (userId: string, file: File) => {
    setUploading(true);
    setUploadError(null);
    try {
      const result = await supabaseStorage.uploadFile(userId, file);
      if (!result) {
        throw new Error('Upload returned empty response');
      }
      toast({
        title: 'Success',
        description: 'Profile picture uploaded successfully.',
      });
      return result;
    } catch (err: any) {
      const errMsg = err.message || 'Failed to upload profile picture.';
      setUploadError(errMsg);
      toast({
        title: 'Error',
        description: errMsg,
        variant: 'destructive',
      });
      return null;
    } finally {
      setUploading(false);
    }
  };

  const remove = async (userId: string) => {
    setDeleting(true);
    setUploadError(null);
    try {
      // Let's implement list and delete in our helper
      const existingFiles = await supabaseStorage.listFiles(userId);
      if (existingFiles.length > 0) {
        const filenames = existingFiles.map(f => f.name);
        const deleted = await supabaseStorage.deleteFiles(userId, filenames);
        if (!deleted) {
          throw new Error('Failed to delete files from storage');
        }
      }
      toast({
        title: 'Success',
        description: 'Profile picture removed successfully.',
      });
      return true;
    } catch (err: any) {
      const errMsg = err.message || 'Failed to remove profile picture.';
      setUploadError(errMsg);
      toast({
        title: 'Error',
        description: errMsg,
        variant: 'destructive',
      });
      return false;
    } finally {
      setDeleting(false);
    }
  };

  return {
    uploading,
    deleting,
    uploadError,
    upload,
    remove,
  };
}
