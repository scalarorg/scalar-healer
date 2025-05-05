package utils


// func BigFloatToNumeric(f *big.Float) (pgtype.Numeric, error) {
// 	str := f.Text('f', -1)
// 	var numeric pgtype.Numeric
// 	err := numeric.Scan(str)
// 	if err != nil {
// 		return pgtype.Numeric{}, err
// 	}
// 	return numeric, nil
// }