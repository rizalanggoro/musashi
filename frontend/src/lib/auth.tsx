import { createContext, useContext, type ReactNode } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { getMeOptions, loginMutation, logoutMutation } from '@/api/@tanstack/react-query.gen'
import type { User } from '@/api/types.gen'

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

  const { data: user, isLoading } = useQuery({
    ...getMeOptions(),
    retry: false,
  })

  const loginMut = useMutation(loginMutation())
  const logoutMut = useMutation(logoutMutation())

  const handleLogin = async (email: string, pass: string) => {
    const res = await loginMut.mutateAsync({
      body: {
        email,
        password: pass,
      },
    })
    if (res?.token) {
      localStorage.setItem('musashi_token', res.token)
    }
    await queryClient.invalidateQueries()
  }

  const handleLogout = async () => {
    try {
      await logoutMut.mutateAsync({})
    } finally {
      localStorage.removeItem('musashi_token')
      queryClient.removeQueries({ queryKey: getMeOptions().queryKey })
      await queryClient.invalidateQueries()
    }
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
        openLoginModal: () => {
          if (typeof window !== 'undefined') {
            window.location.href = '/auth'
          }
        },
      }}
    >
      {children}
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
