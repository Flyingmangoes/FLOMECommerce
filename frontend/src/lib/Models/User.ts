export interface UserModel {
  userId: string;
  firstName: string;
  lastName: string;
  username: string;
  password: string;
  email: string;
  phoneNumber: string;
  locale: string;
  country: string;
  address: string;
  userType: string;
  isVerified: boolean;
  emailConsent: boolean;
  smsConsent: boolean;
  consentUpdated: string;
  consentSrc: string;
  createdAt: string;
  updatedAt: string;
}


export function mapUserResponse(raw: Record<string, unknown>): UserModel {
  return {
    userId: raw.userId as string || '',
    firstName: raw.firstName as string || '',
    lastName: raw.lastName as string || '',
    username: raw.username as string || '',
    password: "",
    email: raw.email as string || '',
    phoneNumber: raw.phoneNumber as string || '',
    locale: raw.locale as string || '',
    country: raw.country as string || '',
    address: raw.address as string || '',
    userType: raw.userType as string || '',  
    isVerified: raw.isVerified as boolean || false,
    emailConsent: raw.emailConsent as boolean || false,
    smsConsent: raw.smsConsent as boolean || false,
    consentUpdated: raw.consentUpdated as string || '',
    consentSrc: raw.consentSrc as string || '',
    createdAt: raw.createdAt as string || '',
    updatedAt: raw.updatedAt as string || ''
    };
}