export interface StaffProfile {
  staff_id: string;
  user_id: string;
  Name: string;
  Username: string;
  Phone?: string | null;
  AvatarURL?: string | null;
  LastLoginAt?: string | null;
  CreatedAt: string;
  UpdatedAt?: string | null;
}
