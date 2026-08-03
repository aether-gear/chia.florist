export interface SupabaseFile {
  name: string;
  id: string;
  updated_at: string;
  created_at: string;
  last_accessed_at: string;
}

export const supabaseStorage = {
  getCredentials() {
    const supabaseUrl = (import.meta.env.SUPABASE_URL || '').replace(/\/$/, '');
    const supabaseKey = import.meta.env.SUPABASE_KEY || '';
    return { supabaseUrl, supabaseKey };
  },

  async listFiles(userId: string): Promise<SupabaseFile[]> {
    const { supabaseUrl, supabaseKey } = this.getCredentials();
    if (!supabaseUrl || !supabaseKey) {
      console.warn('Supabase URL or Key is not configured.');
      return [];
    }

    try {
      const response = await fetch(
        `${supabaseUrl}/storage/v1/object/list/private-assets`,
        {
          method: 'POST',
          headers: {
            'Authorization': `Bearer ${supabaseKey}`,
            'apikey': supabaseKey,
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            prefix: `user/${userId}/`,
            limit: 10,
            sortBy: {
              column: 'name',
              order: 'desc'
            }
          })
        }
      );
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const data = await response.json();
      return data || [];
    } catch (err) {
      console.error('Failed to list files from Supabase:', err);
      return [];
    }
  },

  async deleteFiles(userId: string, filenames: string[]): Promise<boolean> {
    if (filenames.length === 0) return true;
    const { supabaseUrl, supabaseKey } = this.getCredentials();
    if (!supabaseUrl || !supabaseKey) return false;

    try {
      const response = await fetch(
        `${supabaseUrl}/storage/v1/object/private-assets`,
        {
          method: 'DELETE',
          headers: {
            'Authorization': `Bearer ${supabaseKey}`,
            'apikey': supabaseKey,
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            prefixes: filenames.map(name => `user/${userId}/${name}`)
          })
        }
      );
      return response.ok;
    } catch (err) {
      console.error('Failed to delete files from Supabase:', err);
      return false;
    }
  },

  async uploadFile(userId: string, file: File): Promise<{ publicUrl: string; signedUrl: string | null } | null> {
    const { supabaseUrl, supabaseKey } = this.getCredentials();
    if (!supabaseUrl || !supabaseKey) {
      throw new Error('Supabase credentials are not configured in environment.');
    }

    // 1. Delete existing files under user/{userId}/ to ensure only one avatar exists
    const existingFiles = await this.listFiles(userId);
    if (existingFiles.length > 0) {
      const filenames = existingFiles.map(f => f.name);
      await this.deleteFiles(userId, filenames);
    }

    // Sanitize filename to avoid weird character issues
    const sanitizedFilename = file.name.replace(/[^a-zA-Z0-9.-]/g, '_');
    const path = `user/${userId}/${sanitizedFilename}`;

    // 2. Upload new file using POST
    try {
      const response = await fetch(`${supabaseUrl}/storage/v1/object/private-assets/${path}`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${supabaseKey}`,
          'apikey': supabaseKey,
          'Content-Type': file.type
        },
        body: file
      });
      if (!response.ok) {
        const errData = await response.json().catch(() => ({}));
        throw new Error(errData.message || `Upload failed with status: ${response.status}`);
      }
      
      // 3. Directly generate and return signed and public URLs
      const publicUrl = `${supabaseUrl}/storage/v1/object/public/private-assets/${path}`;
      const signedUrl = await this.getSignedUrl(userId, sanitizedFilename);
      return { publicUrl, signedUrl };
    } catch (err: any) {
      console.error('Failed to upload file to Supabase:', err);
      throw new Error(err.message || 'File upload failed', { cause: err });
    }
  },

  async getSignedUrl(userId: string, filename: string): Promise<string | null> {
    const { supabaseUrl, supabaseKey } = this.getCredentials();
    if (!supabaseUrl || !supabaseKey) return null;

    // Extract only the filename just in case a full path was passed
    const cleanFilename = filename.split('/').pop() || filename;
    const path = `user/${userId}/${cleanFilename}`;

    try {
      const response = await fetch(
        `${supabaseUrl}/storage/v1/object/sign/private-assets/${path}`,
        {
          method: 'POST',
          headers: {
            'Authorization': `Bearer ${supabaseKey}`,
            'apikey': supabaseKey,
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            expiresIn: 604800 // 7 days in seconds
          })
        }
      );
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const data = await response.json();
      
      // Supabase returns { signedURL: "..." } or { signedUrl: "..." }
      const signedUrlPath = data?.signedURL || data?.signedUrl || null;
      if (!signedUrlPath) return null;

      // Ensure url is absolute
      let absoluteUrl = signedUrlPath;
      if (!signedUrlPath.startsWith('http') && !signedUrlPath.startsWith('https')) {
        const cleanBase = supabaseUrl.replace(/\/$/, '');
        const cleanPath = signedUrlPath.replace(/^\//, '');
        if (cleanPath.startsWith('storage/v1/')) {
          absoluteUrl = `${cleanBase}/${cleanPath}`;
        } else {
          absoluteUrl = `${cleanBase}/storage/v1/${cleanPath}`;
        }
      }

      return absoluteUrl;
    } catch (err: any) {
      console.error('Failed to generate signed URL:', err);
      return null;
    }
  }
};
