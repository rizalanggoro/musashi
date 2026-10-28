import { createContext, useContext, useState, type ReactNode } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { getMeOptions, loginMutation, logoutMutation } from '@/api/@tanstack/react-query.gen'
import type { User } from '@/api/types.gen'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription } from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { FieldGroup, Field, FieldLabel } from '@/components/ui/field'
import { Alert, AlertDescription } from '@/components/ui/alert'

interface AuthContextValue {
  user: User | undefined
  isLoading: boolean
  isAdmin: boolean
  login: (email: string, pass: string) => Promise<void>
  logout: () => Promise<void>
  openLoginModal: () => void
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined)

export function AuthProvider({ children }: { children: ReactNode }) {
  const queryClient = useQueryClient()
  const [modalOpen, setModalOpen] = useState(false)
  const [email, setEmail] = useState('admin@musashi.local')
  const [password, setPassword] = useState('admin123')
  const [errorMsg, setErrorMsg] = useState('')

  const { data: user, isLoading } = useQuery({
    ...getMeOptions(),
    retry: false,
  })

  const loginMut = useMutation(loginMutation())
  const logoutMut = useMutation(logoutMutation())

  const handleLogin = async (eEmail?: string, ePassword?: string) => {
    setErrorMsg('')
    try {
      await loginMut.mutateAsync({
        body: {
          email: eEmail || email,
          password: ePassword || password,
        },
      })
      await queryClient.invalidateQueries()
      setModalOpen(false)
    } catch (err: unknown) {
      const errorObj = err as { error?: { message?: string }; message?: string }
      const msg = errorObj?.error?.message || errorObj?.message || 'Login gagal: Periksa email dan password.'
      setErrorMsg(msg)
    }
  }

  const handleLogout = async () => {
    await logoutMut.mutateAsync({})
    queryClient.removeQueries({ queryKey: getMeOptions().queryKey })
    await queryClient.invalidateQueries()
  }

  const isAdmin = user?.role === 'admin'

  return (
    <AuthContext.Provider
      value={{
        user,
        isLoading,
        isAdmin,
        login: handleLogin,
        logout: handleLogout,
        openLoginModal: () => setModalOpen(true),
      }}
    >
      {children}

      <Dialog open={modalOpen} onOpenChange={setModalOpen}>
        <DialogContent className="max-w-md">
          <DialogHeader>
            <DialogTitle>Masuk ke MUSASHI</DialogTitle>
            <DialogDescription>
              Gunakan akun default Admin (<code className="text-xs font-semibold">admin@musashi.local</code> / <code className="text-xs font-semibold">admin123</code>).
            </DialogDescription>
          </DialogHeader>

          {errorMsg && (
            <Alert variant="destructive">
              <AlertDescription>{errorMsg}</AlertDescription>
            </Alert>
          )}

          <form
            onSubmit={(e) => {
              e.preventDefault()
              handleLogin()
            }}
            className="flex flex-col gap-4"
          >
            <FieldGroup>
              <Field>
                <FieldLabel htmlFor="auth-email">Email</FieldLabel>
                <Input
                  id="auth-email"
                  type="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  required
                />
              </Field>
              <Field>
                <FieldLabel htmlFor="auth-pass">Password</FieldLabel>
                <Input
                  id="auth-pass"
                  type="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  required
                />
              </Field>
            </FieldGroup>

            <div className="flex justify-end gap-2 pt-2">
              <Button type="button" variant="outline" onClick={() => setModalOpen(false)}>
                Batal
              </Button>
              <Button type="submit" disabled={loginMut.isPending}>
                {loginMut.isPending ? 'Memproses...' : 'Masuk'}
              </Button>
            </div>
          </form>
        </DialogContent>
      </Dialog>
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) {
    throw new Error('useAuth must be used within AuthProvider')
  }
  return ctx
}
