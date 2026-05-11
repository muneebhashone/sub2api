package payment

func YuanToFen(yuanStr string) (int64, error) {
	return AmountToMinorUnit(yuanStr, DefaultPaymentCurrency)
placeholder

func FenToYuan(fen int64) float64 {
	return MinorUnitToAmount(fen, DefaultPaymentCurrency)
placeholder
