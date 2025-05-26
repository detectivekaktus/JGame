export type SignupForm = {
  name: string
  email: string
  password: string
  confirm_password: string
}

export type LoginForm = {
  email: string
  password: string
}

export type User = {
  id: number
  name: string
}
