'use client';

import { useState } from 'react';
import { 
  Save,
  Building,
  DollarSign,
  Shield,
  Users,
  Bell
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Label } from '@/components/ui/label';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import {
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
} from '@/components/ui/tabs';
import { Checkbox } from '@/components/ui/checkbox';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';

// Mock admin users
const adminUsers = [
  { id: '1', name: 'John Admin', email: 'john@greenrideafrica.com', role: 'super_admin', lastLogin: '2024-12-28 14:30', status: 'active' },
  { id: '2', name: 'Sarah Manager', email: 'sarah@greenrideafrica.com', role: 'operations', lastLogin: '2024-12-28 10:15', status: 'active' },
  { id: '3', name: 'Mike Finance', email: 'mike@greenrideafrica.com', role: 'finance', lastLogin: '2024-12-27 16:45', status: 'active' },
  { id: '4', name: 'Jane Support', email: 'jane@greenrideafrica.com', role: 'support', lastLogin: '2024-12-26 09:00', status: 'inactive' },
];

// Mock pricing settings
const pricingSettings = {
  sedan: { baseFare: 1000, perKm: 500, perMin: 100, minFare: 2000 },
  suv: { baseFare: 1500, perKm: 700, perMin: 150, minFare: 3000 },
  moto: { baseFare: 500, perKm: 300, perMin: 50, minFare: 1000 },
  premium: { baseFare: 2500, perKm: 1000, perMin: 200, minFare: 5000 },
};

const getRoleBadge = (role: string) => {
  switch (role) {
    case 'super_admin':
      return <Badge className="bg-purple-100 text-purple-700">Super Admin</Badge>;
    case 'operations':
      return <Badge className="bg-blue-100 text-blue-700">Operations</Badge>;
    case 'finance':
      return <Badge className="bg-green-100 text-green-700">Finance</Badge>;
    case 'support':
      return <Badge className="bg-gray-100 text-gray-700">Support</Badge>;
    default:
      return <Badge variant="secondary">{role}</Badge>;
  }
};

export default function SettingsPage() {
  const [companyName, setCompanyName] = useState('GreenRide Africa');
  const [contactEmail, setContactEmail] = useState('support@greenrideafrica.com');
  const [contactPhone, setContactPhone] = useState('+250 788 000 000');

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div>
        <h1 className="text-2xl font-bold tracking-tight">Settings</h1>
        <p className="text-muted-foreground">
          Manage system configuration and preferences
        </p>
      </div>

      <Tabs defaultValue="general" className="space-y-4">
        <TabsList className="grid w-full grid-cols-5 lg:w-auto lg:grid-cols-none lg:flex">
          <TabsTrigger value="general" className="gap-2">
            <Building className="h-4 w-4" />
            <span className="hidden sm:inline">General</span>
          </TabsTrigger>
          <TabsTrigger value="pricing" className="gap-2">
            <DollarSign className="h-4 w-4" />
            <span className="hidden sm:inline">Pricing</span>
          </TabsTrigger>
          <TabsTrigger value="admins" className="gap-2">
            <Users className="h-4 w-4" />
            <span className="hidden sm:inline">Admin Users</span>
          </TabsTrigger>
          <TabsTrigger value="security" className="gap-2">
            <Shield className="h-4 w-4" />
            <span className="hidden sm:inline">Security</span>
          </TabsTrigger>
          <TabsTrigger value="notifications" className="gap-2">
            <Bell className="h-4 w-4" />
            <span className="hidden sm:inline">Notifications</span>
          </TabsTrigger>
        </TabsList>

        {/* General Settings */}
        <TabsContent value="general">
          <Card>
            <CardHeader>
              <CardTitle>Company Information</CardTitle>
              <CardDescription>
                Basic information about your company
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              <div className="grid gap-4 md:grid-cols-2">
                <div className="space-y-2">
                  <Label htmlFor="company-name">Company Name</Label>
                  <Input
                    id="company-name"
                    value={companyName}
                    onChange={(e) => setCompanyName(e.target.value)}
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="contact-email">Contact Email</Label>
                  <Input
                    id="contact-email"
                    type="email"
                    value={contactEmail}
                    onChange={(e) => setContactEmail(e.target.value)}
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="contact-phone">Contact Phone</Label>
                  <Input
                    id="contact-phone"
                    value={contactPhone}
                    onChange={(e) => setContactPhone(e.target.value)}
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="timezone">Timezone</Label>
                  <Input
                    id="timezone"
                    value="Africa/Kigali (UTC+2)"
                    disabled
                  />
                </div>
              </div>

              <Separator />

              <div className="space-y-4">
                <h4 className="font-medium">Operating Hours</h4>
                <div className="grid gap-4 md:grid-cols-2">
                  <div className="flex items-center justify-between p-3 border rounded-lg">
                    <span>24/7 Service</span>
                    <Checkbox defaultChecked />
                  </div>
                </div>
              </div>

              <div className="flex justify-end">
                <Button className="gap-2">
                  <Save className="h-4 w-4" />
                  Save Changes
                </Button>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        {/* Pricing Settings */}
        <TabsContent value="pricing">
          <Card>
            <CardHeader>
              <CardTitle>Pricing Configuration</CardTitle>
              <CardDescription>
                Set fare rates for different vehicle types
              </CardDescription>
            </CardHeader>
            <CardContent>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Vehicle Type</TableHead>
                    <TableHead className="text-right">Base Fare (RWF)</TableHead>
                    <TableHead className="text-right">Per KM (RWF)</TableHead>
                    <TableHead className="text-right">Per Minute (RWF)</TableHead>
                    <TableHead className="text-right">Minimum Fare (RWF)</TableHead>
                    <TableHead></TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {Object.entries(pricingSettings).map(([type, rates]) => (
                    <TableRow key={type}>
                      <TableCell className="font-medium capitalize">{type}</TableCell>
                      <TableCell className="text-right">{rates.baseFare.toLocaleString()}</TableCell>
                      <TableCell className="text-right">{rates.perKm.toLocaleString()}</TableCell>
                      <TableCell className="text-right">{rates.perMin.toLocaleString()}</TableCell>
                      <TableCell className="text-right">{rates.minFare.toLocaleString()}</TableCell>
                      <TableCell>
                        <Button variant="ghost" size="sm">Edit</Button>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>

              <Separator className="my-6" />

              <div className="space-y-4">
                <h4 className="font-medium">Cancellation Fees</h4>
                <div className="grid gap-4 md:grid-cols-2">
                  <div className="space-y-2">
                    <Label>Passenger Cancellation Fee (RWF)</Label>
                    <Input defaultValue="1000" />
                  </div>
                  <div className="space-y-2">
                    <Label>Free Cancellation Window (minutes)</Label>
                    <Input defaultValue="2" />
                  </div>
                </div>
              </div>

              <div className="flex justify-end mt-6">
                <Button className="gap-2">
                  <Save className="h-4 w-4" />
                  Save Changes
                </Button>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        {/* Admin Users */}
        <TabsContent value="admins">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between">
              <div>
                <CardTitle>Admin Users</CardTitle>
                <CardDescription>
                  Manage dashboard access and permissions
                </CardDescription>
              </div>
              <Button className="gap-2">
                <Users className="h-4 w-4" />
                Add Admin
              </Button>
            </CardHeader>
            <CardContent>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>User</TableHead>
                    <TableHead>Role</TableHead>
                    <TableHead>Last Login</TableHead>
                    <TableHead>Status</TableHead>
                    <TableHead></TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {adminUsers.map((user) => (
                    <TableRow key={user.id}>
                      <TableCell>
                        <div className="flex items-center gap-3">
                          <Avatar className="h-8 w-8">
                            <AvatarFallback className="bg-primary/10 text-primary text-xs">
                              {user.name.split(' ').map(n => n[0]).join('')}
                            </AvatarFallback>
                          </Avatar>
                          <div>
                            <p className="font-medium">{user.name}</p>
                            <p className="text-sm text-muted-foreground">{user.email}</p>
                          </div>
                        </div>
                      </TableCell>
                      <TableCell>{getRoleBadge(user.role)}</TableCell>
                      <TableCell className="text-muted-foreground">{user.lastLogin}</TableCell>
                      <TableCell>
                        {user.status === 'active' ? (
                          <Badge className="bg-green-100 text-green-700">Active</Badge>
                        ) : (
                          <Badge className="bg-gray-100 text-gray-700">Inactive</Badge>
                        )}
                      </TableCell>
                      <TableCell>
                        <Button variant="ghost" size="sm">Edit</Button>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </CardContent>
          </Card>
        </TabsContent>

        {/* Security Settings */}
        <TabsContent value="security">
          <Card>
            <CardHeader>
              <CardTitle>Security Settings</CardTitle>
              <CardDescription>
                Configure security policies and authentication
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              <div className="space-y-4">
                <h4 className="font-medium">Password Policy</h4>
                <div className="space-y-3">
                  <div className="flex items-center justify-between p-3 border rounded-lg">
                    <div>
                      <p className="font-medium">Minimum password length</p>
                      <p className="text-sm text-muted-foreground">At least 8 characters</p>
                    </div>
                    <Input className="w-20" defaultValue="8" />
                  </div>
                  <div className="flex items-center justify-between p-3 border rounded-lg">
                    <div>
                      <p className="font-medium">Require uppercase letters</p>
                      <p className="text-sm text-muted-foreground">At least one uppercase letter</p>
                    </div>
                    <Checkbox defaultChecked />
                  </div>
                  <div className="flex items-center justify-between p-3 border rounded-lg">
                    <div>
                      <p className="font-medium">Require numbers</p>
                      <p className="text-sm text-muted-foreground">At least one number</p>
                    </div>
                    <Checkbox defaultChecked />
                  </div>
                  <div className="flex items-center justify-between p-3 border rounded-lg">
                    <div>
                      <p className="font-medium">Require special characters</p>
                      <p className="text-sm text-muted-foreground">At least one special character</p>
                    </div>
                    <Checkbox />
                  </div>
                </div>
              </div>

              <Separator />

              <div className="space-y-4">
                <h4 className="font-medium">Two-Factor Authentication</h4>
                <div className="flex items-center justify-between p-3 border rounded-lg">
                  <div>
                    <p className="font-medium">Enable 2FA for all admins</p>
                    <p className="text-sm text-muted-foreground">Require two-factor authentication</p>
                  </div>
                  <Checkbox />
                </div>
              </div>

              <Separator />

              <div className="space-y-4">
                <h4 className="font-medium">Session Settings</h4>
                <div className="grid gap-4 md:grid-cols-2">
                  <div className="space-y-2">
                    <Label>Session Timeout (minutes)</Label>
                    <Input defaultValue="60" />
                  </div>
                  <div className="space-y-2">
                    <Label>Max Concurrent Sessions</Label>
                    <Input defaultValue="3" />
                  </div>
                </div>
              </div>

              <div className="flex justify-end">
                <Button className="gap-2">
                  <Save className="h-4 w-4" />
                  Save Changes
                </Button>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        {/* Notification Settings */}
        <TabsContent value="notifications">
          <Card>
            <CardHeader>
              <CardTitle>Notification Preferences</CardTitle>
              <CardDescription>
                Configure email and push notification settings
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              <div className="space-y-4">
                <h4 className="font-medium">Email Notifications</h4>
                <div className="space-y-3">
                  <div className="flex items-center justify-between p-3 border rounded-lg">
                    <div>
                      <p className="font-medium">Daily summary report</p>
                      <p className="text-sm text-muted-foreground">Receive daily stats via email</p>
                    </div>
                    <Checkbox defaultChecked />
                  </div>
                  <div className="flex items-center justify-between p-3 border rounded-lg">
                    <div>
                      <p className="font-medium">New driver registration</p>
                      <p className="text-sm text-muted-foreground">When a new driver signs up</p>
                    </div>
                    <Checkbox defaultChecked />
                  </div>
                  <div className="flex items-center justify-between p-3 border rounded-lg">
                    <div>
                      <p className="font-medium">Payment failures</p>
                      <p className="text-sm text-muted-foreground">When a payment fails</p>
                    </div>
                    <Checkbox defaultChecked />
                  </div>
                  <div className="flex items-center justify-between p-3 border rounded-lg">
                    <div>
                      <p className="font-medium">User complaints</p>
                      <p className="text-sm text-muted-foreground">When a complaint is filed</p>
                    </div>
                    <Checkbox defaultChecked />
                  </div>
                </div>
              </div>

              <div className="flex justify-end">
                <Button className="gap-2">
                  <Save className="h-4 w-4" />
                  Save Changes
                </Button>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  );
}
