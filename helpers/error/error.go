package helpers

func HandleError(err error) error {
	if err != nil {
		return err
	}
	return nil
}
