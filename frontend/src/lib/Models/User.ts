export class UserModels {
    public UserId: string;
    public FirstName: string;
    public LastName: string;
    public Username: string;
    public Email: string;
    public PhoneNumber: string;
    public Locale: string;
    public Country: string;
    public Address: string;
    public UserType: string;
    public IsVerified: boolean;
    public Consent_Updated: string;
    public Consent_src: string;
    public CreatedAt: string;
	public Updatedat: string;

    constructor(
        userId: string,
        username: string, 
        email: string, 
        firstName: string,
        lastName: string,
        phoneNumber: string,
        locale: string,
        country: string,
        address: string,
        userType: string,
        isVerified: boolean,
        emailConsent: boolean,
        smsConsent: boolean,
        consentUpdated: string,
        consentSrc: string,
        createdAt: string,
        updatedAt: string
    ) {
        this.UserId = userId? userId : '';
        this.FirstName = firstName? firstName : '';
        this.LastName = lastName? lastName : '';
        this.Username = username? username : '';
        this.Email = email? email : '';
        this.PhoneNumber = phoneNumber? phoneNumber : '';
        this.Locale = locale? locale : '';
        this.Country = country? country : '';
        this.Address = address? address : '';
        this.UserType = userType? userType : '';
        this.IsVerified = isVerified;
        this.Consent_Updated = consentUpdated? consentUpdated : '';
        this.Consent_src = consentSrc? consentSrc : '';
        this.CreatedAt = createdAt? createdAt : '';
        this.Updatedat = updatedAt? updatedAt : '';
    }
}