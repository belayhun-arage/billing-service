package stripe

type StripeClient struct {}

func (s *StripeClient) Charge(amount int64, currency string) error {
    // call Stripe API
    return nil
}