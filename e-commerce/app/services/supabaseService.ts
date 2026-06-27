// app/services/supabaseService.ts

export interface SupabaseFile {
  name: string
  id: string
  updated_at: string
  created_at: string
  last_accessed_at: string
}

export const supabaseService = {
  getCredentials() {
    const config = useRuntimeConfig()
    const supabaseUrl = ((config.public.supabaseUrl as string) || '').replace(/\/$/, '')
    const supabaseKey = (config.public.supabaseKey as string) || ''
    return { supabaseUrl, supabaseKey }
  },

  async listFiles(userId: string): Promise<SupabaseFile[]> {
    const { supabaseUrl, supabaseKey } = this.getCredentials()
    if (!supabaseUrl || !supabaseKey) {
      console.warn('Supabase URL or Key is not configured.')
      return []
    }

    try {
      const response = await $fetch<SupabaseFile[]>(
        `${supabaseUrl}/storage/v1/object/list/private-assets`,
        {
          method: 'POST',
          headers: {
            'Authorization': `Bearer ${supabaseKey}`,
            'apikey': supabaseKey,
            'Content-Type': 'application/json'
          },
          body: {
            prefix: `user/${userId}/`,
            limit: 10,
            sortBy: {
              column: 'name',
              order: 'desc'
            }
          }
        }
      )
      return response || []
    } catch (err) {
      console.error('Failed to list files from Supabase:', err)
      return []
    }
  },

  async deleteFiles(userId: string, filenames: string[]): Promise<boolean> {
    if (filenames.length === 0) return true
    const { supabaseUrl, supabaseKey } = this.getCredentials()
    if (!supabaseUrl || !supabaseKey) return false

    try {
      await $fetch<any>(
        `${supabaseUrl}/storage/v1/object/private-assets`,
        {
          method: 'DELETE',
          headers: {
            'Authorization': `Bearer ${supabaseKey}`,
            'apikey': supabaseKey,
            'Content-Type': 'application/json'
          },
          body: {
            prefixes: filenames.map(name => `user/${userId}/${name}`)
          }
        }
      )
      return true
    } catch (err) {
      console.error('Failed to delete files from Supabase:', err)
      return false
    }
  },

  async uploadFile(userId: string, file: File): Promise<{ publicUrl: string; signedUrl: string | null } | null> {
    const { supabaseUrl, supabaseKey } = this.getCredentials()
    if (!supabaseUrl || !supabaseKey) {
      throw new Error('Supabase credentials are not configured in environment.')
    }

    // 1. Delete existing files under user/{userId}/ to ensure only one avatar exists
    const existingFiles = await this.listFiles(userId)
    if (existingFiles.length > 0) {
      const filenames = existingFiles.map(f => f.name)
      await this.deleteFiles(userId, filenames)
    }

    // Sanitize filename to avoid weird character issues
    const sanitizedFilename = file.name.replace(/[^a-zA-Z0-9.-]/g, '_')
    const path = `user/${userId}/${sanitizedFilename}`

    // 2. Upload new file using POST
    try {
      await $fetch(`${supabaseUrl}/storage/v1/object/private-assets/${path}`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${supabaseKey}`,
          'apikey': supabaseKey,
          'Content-Type': file.type
        },
        body: file
      })
      
      // 3. Directly generate and return signed and public URLs
      const publicUrl = `${supabaseUrl}/storage/v1/object/public/private-assets/${path}`
      const signedUrl = await this.getSignedUrl(userId, sanitizedFilename)
      return { publicUrl, signedUrl }
    } catch (err: any) {
      console.error('Failed to upload file to Supabase:', err)
      throw new Error(err.data?.message || err.message || 'File upload failed')
    }
  },

  async getSignedUrl(userId: string, filename: string): Promise<string | null> {
    const { supabaseUrl, supabaseKey } = this.getCredentials()
    if (!supabaseUrl || !supabaseKey) return null

    // Extract only the filename just in case a full path was passed
    const cleanFilename = filename.split('/').pop() || filename
    const path = `user/${userId}/${cleanFilename}`

    try {
      const response = await $fetch<any>(
        `${supabaseUrl}/storage/v1/object/sign/private-assets/${path}`,
        {
          method: 'POST',
          headers: {
            'Authorization': `Bearer ${supabaseKey}`,
            'apikey': supabaseKey,
            'Content-Type': 'application/json'
          },
          body: {
            expiresIn: 604800 // 7 days in seconds
          }
        }
      )
      // Supabase returns { signedURL: "..." } or { signedUrl: "..." }
      var signedUrl = response?.signedURL || response?.signedUrl || null
      
      // Ensure url is absolute if relative
      signedUrl = "storage/v1/" + signedUrl
      let absoluteUrl = signedUrl
      if (signedUrl && !signedUrl.startsWith('http') && !signedUrl.startsWith('https')) {
        const cleanBase = supabaseUrl.replace(/\/$/, '')
        const cleanPath = signedUrl.replace(/^\//, '')
        absoluteUrl = `${cleanBase}/${cleanPath}`
      }

      return absoluteUrl
    } catch (err: any) {
      console.error('Failed to generate signed URL:', err)
      if (err.data) {
        console.error('Supabase signed URL error payload:', err.data)
      }
      return null
    }
  },

  async getAvatarUrls(userId: string): Promise<{ publicUrl: string; signedUrl: string | null } | null> {
    const { supabaseUrl } = this.getCredentials()
    if (!supabaseUrl) return null

    const files = await this.listFiles(userId)
    if (files.length === 0) return null

    const rawFilename = files[0]!.name
    const filename = rawFilename.split('/').pop() || rawFilename
    const path = `user/${userId}/${filename}`
    
    // Construct standard public URL as fallback
    const publicUrl = `${supabaseUrl}/storage/v1/object/public/private-assets/${path}`
    
    // Fetch signed URL in case it's a private bucket
    const signedUrl = await this.getSignedUrl(userId, filename)

    return {
      publicUrl,
      signedUrl
    }
  }
}
