export interface SignupForm {
  name:             string
  email:            string
  password:         string
  confirm_password: string
}

export interface LoginForm {
  email:    string
  password: string
}

export interface User {
  id:             number
  name:           string
  matches_played: number
  matches_won:    number
}

export interface Me extends User {
  email: string
}

