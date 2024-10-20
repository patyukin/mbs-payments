package model

func (lp SignUpData) IsValid() bool {
	if lp.Login == "" || lp.Password == "" || lp.Telegram == "" {
		return false
	}

	return true
}
