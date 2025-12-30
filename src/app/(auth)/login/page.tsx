'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { Leaf, Mail, Lock, Eye, EyeOff, AlertCircle } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Checkbox } from '@/components/ui/checkbox';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { useAuthStore } from '@/stores/auth-store';
import { cn } from '@/lib/utils';

export default function LoginPage() {
  const router = useRouter();
  const { setUser } = useAuthStore();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [rememberMe, setRememberMe] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setIsLoading(true);

    try {
      // For demo purposes, accept any credentials
      // In production, this would call the API
      await new Promise((resolve) => setTimeout(resolve, 1000));

      if (email && password) {
        // Mock successful login
        setUser({
          id: '1',
          name: 'Admin User',
          email: email,
          role: 'super_admin',
          permissions: ['read', 'write', 'delete', 'export', 'settings'],
          createdAt: new Date(),
        });

        // Store mock token
        if (typeof window !== 'undefined') {
          localStorage.setItem('admin_token', 'mock_token_123');
        }

        router.push('/');
      } else {
        setError('Please enter email and password');
      }
    } catch (err) {
      setError('Invalid email or password');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex">
      {/* Left Side - Branding */}
      <div className="hidden lg:flex lg:w-1/2 bg-gradient-to-br from-primary to-green-600 p-12 flex-col justify-between">
        <div>
          <div className="flex items-center gap-3">
            <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-white/20 backdrop-blur">
              <Leaf className="h-7 w-7 text-white" />
            </div>
            <span className="text-2xl font-bold text-white">GreenRide</span>
          </div>
          <p className="mt-2 text-white/80">Admin Dashboard</p>
        </div>

        <div className="space-y-6">
          <h1 className="text-4xl font-bold text-white leading-tight">
            Manage your ride-hailing platform with ease
          </h1>
          <p className="text-lg text-white/80">
            Monitor drivers, track revenue, and analyze performance all in one place.
          </p>
          <div className="grid grid-cols-3 gap-4">
            <div className="rounded-xl bg-white/10 backdrop-blur p-4">
              <p className="text-3xl font-bold text-white">500+</p>
              <p className="text-sm text-white/70">Active Drivers</p>
            </div>
            <div className="rounded-xl bg-white/10 backdrop-blur p-4">
              <p className="text-3xl font-bold text-white">10K+</p>
              <p className="text-sm text-white/70">Daily Rides</p>
            </div>
            <div className="rounded-xl bg-white/10 backdrop-blur p-4">
              <p className="text-3xl font-bold text-white">50K+</p>
              <p className="text-sm text-white/70">Happy Users</p>
            </div>
          </div>
        </div>

        <p className="text-white/60 text-sm">
          © 2025 GreenRide Africa. All rights reserved.
        </p>
      </div>

      {/* Right Side - Login Form */}
      <div className="flex-1 flex items-center justify-center p-8 bg-background">
        <div className="w-full max-w-md">
          {/* Mobile Logo */}
          <div className="lg:hidden flex items-center justify-center gap-2 mb-8">
            <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary">
              <Leaf className="h-6 w-6 text-primary-foreground" />
            </div>
            <span className="text-xl font-bold">GreenRide Admin</span>
          </div>

          <Card className="border-0 shadow-lg">
            <CardHeader className="space-y-1 text-center">
              <CardTitle className="text-2xl font-bold">Welcome back</CardTitle>
              <CardDescription>
                Sign in to your admin account
              </CardDescription>
            </CardHeader>
            <CardContent>
              <form onSubmit={handleSubmit} className="space-y-4">
                {error && (
                  <div className="flex items-center gap-2 rounded-lg bg-destructive/10 p-3 text-sm text-destructive">
                    <AlertCircle className="h-4 w-4" />
                    {error}
                  </div>
                )}

                <div className="space-y-2">
                  <Label htmlFor="email">Email address</Label>
                  <div className="relative">
                    <Mail className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
                    <Input
                      id="email"
                      type="email"
                      placeholder="admin@greenrideafrica.com"
                      className="pl-10"
                      value={email}
                      onChange={(e) => setEmail(e.target.value)}
                      disabled={isLoading}
                      required
                    />
                  </div>
                </div>

                <div className="space-y-2">
                  <div className="flex items-center justify-between">
                    <Label htmlFor="password">Password</Label>
                    <Button
                      type="button"
                      variant="link"
                      className="h-auto p-0 text-xs text-primary"
                    >
                      Forgot password?
                    </Button>
                  </div>
                  <div className="relative">
                    <Lock className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
                    <Input
                      id="password"
                      type={showPassword ? 'text' : 'password'}
                      placeholder="••••••••"
                      className="pl-10 pr-10"
                      value={password}
                      onChange={(e) => setPassword(e.target.value)}
                      disabled={isLoading}
                      required
                    />
                    <Button
                      type="button"
                      variant="ghost"
                      size="icon"
                      className="absolute right-1 top-1/2 h-7 w-7 -translate-y-1/2"
                      onClick={() => setShowPassword(!showPassword)}
                    >
                      {showPassword ? (
                        <EyeOff className="h-4 w-4 text-muted-foreground" />
                      ) : (
                        <Eye className="h-4 w-4 text-muted-foreground" />
                      )}
                    </Button>
                  </div>
                </div>

                <div className="flex items-center space-x-2">
                  <Checkbox
                    id="remember"
                    checked={rememberMe}
                    onCheckedChange={(checked) => setRememberMe(checked as boolean)}
                  />
                  <Label htmlFor="remember" className="text-sm font-normal">
                    Remember me
                  </Label>
                </div>

                <Button
                  type="submit"
                  className="w-full"
                  disabled={isLoading}
                >
                  {isLoading ? (
                    <div className="flex items-center gap-2">
                      <div className="h-4 w-4 animate-spin rounded-full border-2 border-current border-t-transparent" />
                      Signing in...
                    </div>
                  ) : (
                    'Sign In'
                  )}
                </Button>
              </form>

              <div className="mt-6 text-center">
                <p className="text-sm text-muted-foreground">
                  Need help?{' '}
                  <Button variant="link" className="h-auto p-0 text-primary">
                    Contact support
                  </Button>
                </p>
              </div>
            </CardContent>
          </Card>

          {/* Demo credentials hint */}
          <p className="mt-4 text-center text-xs text-muted-foreground">
            Demo: Use any email and password to login
          </p>
        </div>
      </div>
    </div>
  );
}
