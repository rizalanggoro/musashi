import { useState } from 'react'
import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useAuth } from '@/lib/auth'
import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from '@/components/ui/card'
import { FieldGroup, Field, FieldLabel } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { LogIn, Construction } from 'lucide-react'

export const Route = createFileRoute('/auth')({
  component: AuthPage,
})

function AuthPage() {
  const { login } = useAuth()
  const navigate = useNavigate()
  const [email, setEmail] = useState('admin@musashi.local')
  const [password, setPassword] = useState('admin123')
  const [errorMsg, setErrorMsg] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setErrorMsg('')
    setIsSubmitting(true)
    try {
      await login(email, password)
      navigate({ to: '/projects', replace: true })
    } catch (err: unknown) {
      const errorObj = err as { error?: { message?: string }; message?: string }
      const msg = errorObj?.error?.message || errorObj?.message || 'Login gagal: Periksa email dan password.'
      setErrorMsg(msg)
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <div className="flex min-h-screen w-full items-center justify-center bg-muted/40 p-4">
      <Card className="w-full max-w-sm">
        <CardHeader className="text-center">
          <div className="mx-auto mb-2 flex size-10 items-center justify-center rounded-lg bg-primary text-primary-foreground">
            <Construction className="size-5" />
          </div>
          <CardTitle className="text-xl">MUSASHI</CardTitle>
          <CardDescription>
            Estimasi Biaya &amp; Manajemen Proyek Konstruksi
          </CardDescription>
        </CardHeader>

        <CardContent>
          <form onSubmit={handleSubmit} className="flex flex-col gap-4">
            {errorMsg && (
              <Alert variant="destructive">
                <AlertDescription>{errorMsg}</AlertDescription>
              </Alert>
            )}

            <FieldGroup>
              <Field>
                <FieldLabel htmlFor="auth-email">Email</FieldLabel>
                <Input
                  id="auth-email"
                  type="email"
                  placeholder="admin@musashi.local"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  autoComplete="email"
                  required
                />
              </Field>

              <Field>
                <FieldLabel htmlFor="auth-password">Password</FieldLabel>
                <Input
                  id="auth-password"
                  type="password"
                  placeholder="••••••••"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  autoComplete="current-password"
                  required
                />
              </Field>
            </FieldGroup>

            <Button type="submit" className="w-full" disabled={isSubmitting}>
              <LogIn data-icon="inline-start" />
              {isSubmitting ? 'Memproses...' : 'Masuk'}
            </Button>
          </form>
        </CardContent>

        <CardFooter className="flex flex-col gap-1 border-t pt-4 text-center text-xs text-muted-foreground">
          <span>Akun default admin:</span>
          <span className="font-mono text-[11px] text-foreground">
            admin@musashi.local / admin123
          </span>
        </CardFooter>
      </Card>
    </div>
  )
}

