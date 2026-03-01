'use client';

import { useState, useEffect } from 'react';
import {
  Save,
  Building,
  DollarSign,
  Shield,
  Users,
  Bell,
  Loader2,
  CheckCircle,
  AlertTriangle,
  Wrench,
  Phone,
  Database,
  Trash2,
  Smartphone
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
import { Textarea } from '@/components/ui/textarea';
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import { toast } from 'sonner';
import { apiClient } from '@/lib/api-client';

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
  // Support config state
  const [isLoadingConfig, setIsLoadingConfig] = useState(true);
  const [isSavingGeneral, setIsSavingGeneral] = useState(false);
  const [isSavingPricing, setIsSavingPricing] = useState(false);
  const [isSavingSecurity, setIsSavingSecurity] = useState(false);
  const [isSavingNotifications, setIsSavingNotifications] = useState(false);
  
  const [companyName, setCompanyName] = useState('GreenRide Africa');
  const [contactEmail, setContactEmail] = useState('');
  const [contactPhone, setContactPhone] = useState('');
  const [contactWhatsApp, setContactWhatsApp] = useState('');
  const [operatingHours, setOperatingHours] = useState('24/7');
  const [is24x7, setIs24x7] = useState(true);

  // System config state
  const [maintenanceMode, setMaintenanceMode] = useState(false);
  const [maintenanceMessage, setMaintenanceMessage] = useState('');
  const [maintenancePhone, setMaintenancePhone] = useState('6996');
  const [maintenanceStartedAt, setMaintenanceStartedAt] = useState(0);
  const [updateNoticeEnabled, setUpdateNoticeEnabled] = useState(false);
  const [forceUpdateEnabled, setForceUpdateEnabled] = useState(false);
  const [updateNoticeTitle, setUpdateNoticeTitle] = useState('Update available');
  const [updateNoticeMessage, setUpdateNoticeMessage] = useState('');
  const [minimumAppVersion, setMinimumAppVersion] = useState('');
  const [latestAppVersion, setLatestAppVersion] = useState('');
  const [androidStoreUrl, setAndroidStoreUrl] = useState('');
  const [iosStoreUrl, setIosStoreUrl] = useState('');
  const [isLoadingSystem, setIsLoadingSystem] = useState(true);
  const [isSavingSystem, setIsSavingSystem] = useState(false);
  const [showMaintenanceConfirm, setShowMaintenanceConfirm] = useState(false);
  const [confirmText, setConfirmText] = useState('');
  const [showPurgeConfirm, setShowPurgeConfirm] = useState(false);
  const [purgeConfirmText, setPurgeConfirmText] = useState('');
  const [purgeDryRun, setPurgeDryRun] = useState(true);
  const [isPurging, setIsPurging] = useState(false);
  const [purgeResult, setPurgeResult] = useState<{
    dry_run: boolean;
    zombie_users_found: number;
    zombie_users_deleted: number;
    vehicle_unbound_count: number;
    fcm_tokens_deleted_count: number;
    vehicles_deleted_count: number;
    price_rules_deleted: number;
    promotions_deleted: number;
    ratings_deleted: number;
  } | null>(null);

  // Load support config on mount
  useEffect(() => {
    const loadConfig = async () => {
      try {
        const response = await apiClient.getSupportConfig();
        if (response.data) {
          setContactEmail(response.data.email || '');
          setContactPhone(response.data.phone || '');
          setContactWhatsApp(response.data.whatsapp || '');
          setOperatingHours(response.data.hours || '24/7');
          setIs24x7(response.data.hours === '24/7');
        }
      } catch (error) {
        console.error('Failed to load support config:', error);
        // Use defaults on error
        setContactEmail('support@greenrideafrica.com');
        setContactPhone('+250 788 000 000');
      } finally {
        setIsLoadingConfig(false);
      }
    };
    loadConfig();
  }, []);

  // Load system config on mount
  useEffect(() => {
    const loadSystemConfig = async () => {
      try {
        const response = await apiClient.getSystemConfig();
        if (response.data) {
          setMaintenanceMode(response.data.maintenance_mode);
          setMaintenanceMessage(response.data.maintenance_message || '');
          setMaintenancePhone(response.data.maintenance_phone || '6996');
          setMaintenanceStartedAt(response.data.maintenance_started_at || 0);
          setUpdateNoticeEnabled(response.data.update_notice_enabled || false);
          setForceUpdateEnabled(response.data.force_update_enabled || false);
          setUpdateNoticeTitle(response.data.update_notice_title || 'Update available');
          setUpdateNoticeMessage(response.data.update_notice_message || '');
          setMinimumAppVersion(response.data.minimum_app_version || '');
          setLatestAppVersion(response.data.latest_app_version || '');
          setAndroidStoreUrl(response.data.android_store_url || '');
          setIosStoreUrl(response.data.ios_store_url || '');
        }
      } catch (error) {
        console.error('Failed to load system config:', error);
      } finally {
        setIsLoadingSystem(false);
      }
    };
    loadSystemConfig();
  }, []);

  // Save general settings
  const handleSaveGeneral = async () => {
    setIsSavingGeneral(true);
    try {
      await apiClient.updateSupportConfig({
        email: contactEmail,
        phone: contactPhone,
        whatsapp: contactWhatsApp,
        hours: is24x7 ? '24/7' : operatingHours,
      });
      toast.success('General settings saved successfully!', {
        icon: <CheckCircle className="h-4 w-4 text-green-500" />,
      });
    } catch {
      toast.error('Failed to save settings. Please try again.');
    } finally {
      setIsSavingGeneral(false);
    }
  };

  // Save pricing settings (mock)
  const handleSavePricing = async () => {
    setIsSavingPricing(true);
    // Simulate API call
    await new Promise(r => setTimeout(r, 500));
    setIsSavingPricing(false);
    toast.success('Pricing settings saved successfully!', {
      icon: <CheckCircle className="h-4 w-4 text-green-500" />,
    });
  };

  // Save security settings (mock)
  const handleSaveSecurity = async () => {
    setIsSavingSecurity(true);
    await new Promise(r => setTimeout(r, 500));
    setIsSavingSecurity(false);
    toast.success('Security settings saved successfully!', {
      icon: <CheckCircle className="h-4 w-4 text-green-500" />,
    });
  };

  // Handle maintenance mode toggle
  const handleMaintenanceToggle = () => {
    if (!maintenanceMode) {
      // Enabling maintenance - show confirmation dialog
      setConfirmText('');
      setShowMaintenanceConfirm(true);
    } else {
      // Disabling maintenance - save immediately
      handleSaveSystem(false);
    }
  };

  // Confirm enable maintenance mode
  const handleConfirmMaintenance = () => {
    setShowMaintenanceConfirm(false);
    setConfirmText('');
    handleSaveSystem(true);
  };

  // Save system config
  const handleSaveSystem = async (newMaintenanceMode: boolean) => {
    setIsSavingSystem(true);
    try {
      const response = await apiClient.updateSystemConfig({
        maintenance_mode: newMaintenanceMode,
        maintenance_message: maintenanceMessage,
        maintenance_phone: maintenancePhone,
        update_notice_enabled: updateNoticeEnabled,
        force_update_enabled: forceUpdateEnabled,
        update_notice_title: updateNoticeTitle,
        update_notice_message: updateNoticeMessage,
        minimum_app_version: minimumAppVersion,
        latest_app_version: latestAppVersion,
        android_store_url: androidStoreUrl,
        ios_store_url: iosStoreUrl,
      });
      if (response.data) {
        setMaintenanceMode(response.data.maintenance_mode);
        setMaintenanceStartedAt(response.data.maintenance_started_at || 0);
        setUpdateNoticeEnabled(response.data.update_notice_enabled || false);
        setForceUpdateEnabled(response.data.force_update_enabled || false);
        setUpdateNoticeTitle(response.data.update_notice_title || 'Update available');
        setUpdateNoticeMessage(response.data.update_notice_message || '');
        setMinimumAppVersion(response.data.minimum_app_version || '');
        setLatestAppVersion(response.data.latest_app_version || '');
        setAndroidStoreUrl(response.data.android_store_url || '');
        setIosStoreUrl(response.data.ios_store_url || '');
      }
      toast.success(
        newMaintenanceMode
          ? 'Maintenance mode ENABLED. Users will see the maintenance screen.'
          : 'Maintenance mode disabled. Service is back to normal.',
        {
          icon: newMaintenanceMode
            ? <AlertTriangle className="h-4 w-4 text-amber-500" />
            : <CheckCircle className="h-4 w-4 text-green-500" />,
        }
      );
    } catch {
      toast.error('Failed to update system config. Please try again.');
    } finally {
      setIsSavingSystem(false);
    }
  };

  // Save maintenance message/phone without toggling mode
  const handleSaveSystemSettings = async () => {
    setIsSavingSystem(true);
    try {
      const response = await apiClient.updateSystemConfig({
        maintenance_mode: maintenanceMode,
        maintenance_message: maintenanceMessage,
        maintenance_phone: maintenancePhone,
        update_notice_enabled: updateNoticeEnabled,
        force_update_enabled: forceUpdateEnabled,
        update_notice_title: updateNoticeTitle,
        update_notice_message: updateNoticeMessage,
        minimum_app_version: minimumAppVersion,
        latest_app_version: latestAppVersion,
        android_store_url: androidStoreUrl,
        ios_store_url: iosStoreUrl,
      });
      if (response.data) {
        setMaintenanceMode(response.data.maintenance_mode);
        setMaintenanceStartedAt(response.data.maintenance_started_at || 0);
        setUpdateNoticeEnabled(response.data.update_notice_enabled || false);
        setForceUpdateEnabled(response.data.force_update_enabled || false);
        setUpdateNoticeTitle(response.data.update_notice_title || 'Update available');
        setUpdateNoticeMessage(response.data.update_notice_message || '');
        setMinimumAppVersion(response.data.minimum_app_version || '');
        setLatestAppVersion(response.data.latest_app_version || '');
        setAndroidStoreUrl(response.data.android_store_url || '');
        setIosStoreUrl(response.data.ios_store_url || '');
      }
      toast.success('System settings saved successfully!', {
        icon: <CheckCircle className="h-4 w-4 text-green-500" />,
      });
    } catch {
      toast.error('Failed to save system settings. Please try again.');
    } finally {
      setIsSavingSystem(false);
    }
  };

  const handleOpenPurgeDialog = (dryRun: boolean) => {
    setPurgeDryRun(dryRun);
    setPurgeConfirmText('');
    setShowPurgeConfirm(true);
  };

  const handleRunPurgeCleanup = async () => {
    setIsPurging(true);
    try {
      const response = await apiClient.purgeLegacyDeleted({
        confirm: 'PURGE_LEGACY_DELETED',
        dry_run: purgeDryRun,
      });
      setPurgeResult(response.data);
      toast.success(
        purgeDryRun
          ? 'Dry-run completed successfully.'
          : 'Legacy deleted cleanup completed.',
        { icon: <CheckCircle className="h-4 w-4 text-green-500" /> }
      );
      setShowPurgeConfirm(false);
      setPurgeConfirmText('');
    } catch {
      toast.error('Failed to run cleanup. Please try again.');
    } finally {
      setIsPurging(false);
    }
  };

  // Save notification settings (mock)
  const handleSaveNotifications = async () => {
    setIsSavingNotifications(true);
    await new Promise(r => setTimeout(r, 500));
    setIsSavingNotifications(false);
    toast.success('Notification preferences saved successfully!', {
      icon: <CheckCircle className="h-4 w-4 text-green-500" />,
    });
  };

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
        <TabsList className="grid w-full grid-cols-6 lg:w-auto lg:grid-cols-none lg:flex">
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
          <TabsTrigger value="system" className="gap-2">
            <Wrench className="h-4 w-4" />
            <span className="hidden sm:inline">System</span>
          </TabsTrigger>
        </TabsList>

        {/* General Settings */}
        <TabsContent value="general">
          <Card>
            <CardHeader>
              <CardTitle>Company Information</CardTitle>
              <CardDescription>
                Basic information about your company and support contact details
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              {isLoadingConfig ? (
                <div className="flex items-center justify-center py-8">
                  <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
                  <span className="ml-2 text-muted-foreground">Loading settings...</span>
                </div>
              ) : (
                <>
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
                      <Label htmlFor="contact-email">Support Email</Label>
                      <Input
                        id="contact-email"
                        type="email"
                        value={contactEmail}
                        onChange={(e) => setContactEmail(e.target.value)}
                        placeholder="support@greenrideafrica.com"
                      />
                    </div>
                    <div className="space-y-2">
                      <Label htmlFor="contact-phone">Support Phone</Label>
                      <Input
                        id="contact-phone"
                        value={contactPhone}
                        onChange={(e) => setContactPhone(e.target.value)}
                        placeholder="+250 788 000 000"
                      />
                    </div>
                    <div className="space-y-2">
                      <Label htmlFor="contact-whatsapp">WhatsApp Number</Label>
                      <Input
                        id="contact-whatsapp"
                        value={contactWhatsApp}
                        onChange={(e) => setContactWhatsApp(e.target.value)}
                        placeholder="+250 788 000 001"
                      />
                      <p className="text-xs text-muted-foreground">
                        Used for WhatsApp support in mobile app
                      </p>
                    </div>
                  </div>

                  <Separator />

                  <div className="grid gap-4 md:grid-cols-2">
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
                        <Checkbox 
                          checked={is24x7} 
                          onCheckedChange={(checked) => setIs24x7(checked === true)}
                        />
                      </div>
                    </div>
                  </div>

                  <div className="flex justify-end">
                    <Button 
                      className="gap-2" 
                      onClick={handleSaveGeneral}
                      disabled={isSavingGeneral}
                    >
                      {isSavingGeneral ? (
                        <Loader2 className="h-4 w-4 animate-spin" />
                      ) : (
                        <Save className="h-4 w-4" />
                      )}
                      {isSavingGeneral ? 'Saving...' : 'Save Changes'}
                    </Button>
                  </div>
                </>
              )}
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
                <Button 
                  className="gap-2" 
                  onClick={handleSavePricing}
                  disabled={isSavingPricing}
                >
                  {isSavingPricing ? (
                    <Loader2 className="h-4 w-4 animate-spin" />
                  ) : (
                    <Save className="h-4 w-4" />
                  )}
                  {isSavingPricing ? 'Saving...' : 'Save Changes'}
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
                <Button 
                  className="gap-2" 
                  onClick={handleSaveSecurity}
                  disabled={isSavingSecurity}
                >
                  {isSavingSecurity ? (
                    <Loader2 className="h-4 w-4 animate-spin" />
                  ) : (
                    <Save className="h-4 w-4" />
                  )}
                  {isSavingSecurity ? 'Saving...' : 'Save Changes'}
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
                <Button 
                  className="gap-2" 
                  onClick={handleSaveNotifications}
                  disabled={isSavingNotifications}
                >
                  {isSavingNotifications ? (
                    <Loader2 className="h-4 w-4 animate-spin" />
                  ) : (
                    <Save className="h-4 w-4" />
                  )}
                  {isSavingNotifications ? 'Saving...' : 'Save Changes'}
                </Button>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        {/* System Settings (Maintenance Mode) */}
        <TabsContent value="system">
          <div className="space-y-4">
            {/* Maintenance Mode Status Banner */}
            {maintenanceMode && (
              <div className="flex items-center gap-3 p-4 bg-amber-50 border border-amber-200 rounded-lg">
                <AlertTriangle className="h-5 w-5 text-amber-600 shrink-0" />
                <div className="flex-1">
                  <p className="font-medium text-amber-800">Maintenance Mode is ACTIVE</p>
                  <p className="text-sm text-amber-600">
                    Users are currently seeing the maintenance screen.
                    {maintenanceStartedAt > 0 && (
                      <> Started {new Date(maintenanceStartedAt * 1000).toLocaleString()}.</>
                    )}
                  </p>
                </div>
                <Button
                  variant="outline"
                  size="sm"
                  className="border-amber-300 text-amber-700 hover:bg-amber-100"
                  onClick={handleMaintenanceToggle}
                  disabled={isSavingSystem}
                >
                  {isSavingSystem ? <Loader2 className="h-4 w-4 animate-spin" /> : 'Disable'}
                </Button>
              </div>
            )}

            <Card>
              <CardHeader>
                <CardTitle>Maintenance Mode</CardTitle>
                <CardDescription>
                  When enabled, all mobile app users will see a maintenance screen and cannot use the service.
                  Active rides (in-progress) will NOT be interrupted.
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-6">
                {isLoadingSystem ? (
                  <div className="flex items-center justify-center py-8">
                    <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
                    <span className="ml-2 text-muted-foreground">Loading system config...</span>
                  </div>
                ) : (
                  <>
                    {/* Toggle */}
                    <div
                      className={`flex items-center justify-between p-4 border-2 rounded-lg cursor-pointer transition-colors ${
                        maintenanceMode
                          ? 'border-red-300 bg-red-50'
                          : 'border-gray-200 hover:border-gray-300'
                      }`}
                      onClick={handleMaintenanceToggle}
                    >
                      <div className="flex items-center gap-3">
                        <div className={`p-2 rounded-full ${maintenanceMode ? 'bg-red-100' : 'bg-gray-100'}`}>
                          <Wrench className={`h-5 w-5 ${maintenanceMode ? 'text-red-600' : 'text-gray-500'}`} />
                        </div>
                        <div>
                          <p className="font-medium">
                            {maintenanceMode ? 'Maintenance Mode is ON' : 'Maintenance Mode is OFF'}
                          </p>
                          <p className="text-sm text-muted-foreground">
                            {maintenanceMode
                              ? 'Users cannot access the app. Click to disable.'
                              : 'Service is running normally. Click to enable maintenance.'}
                          </p>
                        </div>
                      </div>
                      <div
                        className={`w-12 h-7 rounded-full transition-colors relative ${
                          maintenanceMode ? 'bg-red-500' : 'bg-gray-300'
                        }`}
                      >
                        <div
                          className={`absolute top-0.5 w-6 h-6 rounded-full bg-white shadow transition-transform ${
                            maintenanceMode ? 'translate-x-5' : 'translate-x-0.5'
                          }`}
                        />
                      </div>
                    </div>

                    <Separator />

                    {/* Maintenance Message */}
                    <div className="space-y-2">
                      <Label htmlFor="maintenance-message">Maintenance Message</Label>
                      <Textarea
                        id="maintenance-message"
                        value={maintenanceMessage}
                        onChange={(e) => setMaintenanceMessage(e.target.value)}
                        placeholder="We're currently improving your experience. Our services will resume shortly."
                        rows={3}
                      />
                      <p className="text-xs text-muted-foreground">
                        This message is shown to users on the maintenance screen in the mobile app.
                      </p>
                    </div>

                    {/* Support Phone */}
                    <div className="space-y-2">
                      <Label htmlFor="maintenance-phone">
                        <span className="flex items-center gap-2">
                          <Phone className="h-4 w-4" />
                          Emergency Support Phone
                        </span>
                      </Label>
                      <Input
                        id="maintenance-phone"
                        value={maintenancePhone}
                        onChange={(e) => setMaintenancePhone(e.target.value)}
                        placeholder="6996"
                      />
                      <p className="text-xs text-muted-foreground">
                        Shown on the maintenance screen so users can call for urgent matters.
                      </p>
                    </div>

                    <div className="flex justify-end">
                      <Button
                        className="gap-2"
                        onClick={handleSaveSystemSettings}
                        disabled={isSavingSystem}
                      >
                        {isSavingSystem ? (
                          <Loader2 className="h-4 w-4 animate-spin" />
                        ) : (
                          <Save className="h-4 w-4" />
                        )}
                        {isSavingSystem ? 'Saving...' : 'Save Settings'}
                      </Button>
                    </div>
                  </>
                )}
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Smartphone className="h-5 w-5" />
                  App Update Notice
                </CardTitle>
                <CardDescription>
                  Control in-app update notice payload for mobile clients via <code>/system/config</code>.
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-6">
                <div
                  className={`flex items-center justify-between p-4 border-2 rounded-lg cursor-pointer transition-colors ${
                    updateNoticeEnabled
                      ? 'border-blue-300 bg-blue-50'
                      : 'border-gray-200 hover:border-gray-300'
                  }`}
                  onClick={() => setUpdateNoticeEnabled(!updateNoticeEnabled)}
                >
                  <div>
                    <p className="font-medium">
                      {updateNoticeEnabled ? 'Update Notice is ON' : 'Update Notice is OFF'}
                    </p>
                    <p className="text-sm text-muted-foreground">
                      {updateNoticeEnabled
                        ? 'Mobile app can show an update prompt based on this config.'
                        : 'Enable to send update notice config to mobile app.'}
                    </p>
                  </div>
                  <div
                    className={`w-12 h-7 rounded-full transition-colors relative ${
                      updateNoticeEnabled ? 'bg-blue-500' : 'bg-gray-300'
                    }`}
                  >
                    <div
                      className={`absolute top-0.5 w-6 h-6 rounded-full bg-white shadow transition-transform ${
                        updateNoticeEnabled ? 'translate-x-5' : 'translate-x-0.5'
                      }`}
                    />
                  </div>
                </div>

                <div className="flex items-center justify-between p-3 border rounded-lg">
                  <div>
                    <p className="font-medium">Force Update</p>
                    <p className="text-sm text-muted-foreground">
                      If enabled, older app versions should be blocked by the mobile client.
                    </p>
                  </div>
                  <Checkbox
                    checked={forceUpdateEnabled}
                    onCheckedChange={(checked) => setForceUpdateEnabled(checked === true)}
                  />
                </div>

                <div className="grid gap-4 md:grid-cols-2">
                  <div className="space-y-2">
                    <Label htmlFor="minimum-app-version">Minimum App Version</Label>
                    <Input
                      id="minimum-app-version"
                      value={minimumAppVersion}
                      onChange={(e) => setMinimumAppVersion(e.target.value)}
                      placeholder="e.g. 1.5.6+84"
                    />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="latest-app-version">Latest App Version</Label>
                    <Input
                      id="latest-app-version"
                      value={latestAppVersion}
                      onChange={(e) => setLatestAppVersion(e.target.value)}
                      placeholder="e.g. 1.5.7+85"
                    />
                  </div>
                </div>

                <div className="space-y-2">
                  <Label htmlFor="update-notice-title">Notice Title</Label>
                  <Input
                    id="update-notice-title"
                    value={updateNoticeTitle}
                    onChange={(e) => setUpdateNoticeTitle(e.target.value)}
                    placeholder="Update available"
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="update-notice-message">Notice Message</Label>
                  <Textarea
                    id="update-notice-message"
                    value={updateNoticeMessage}
                    onChange={(e) => setUpdateNoticeMessage(e.target.value)}
                    placeholder="A new version of Green Ride is available. Please update for the best experience."
                    rows={3}
                  />
                </div>

                <div className="grid gap-4 md:grid-cols-2">
                  <div className="space-y-2">
                    <Label htmlFor="android-store-url">Android Store URL</Label>
                    <Input
                      id="android-store-url"
                      value={androidStoreUrl}
                      onChange={(e) => setAndroidStoreUrl(e.target.value)}
                      placeholder="https://play.google.com/store/apps/details?id=..."
                    />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="ios-store-url">iOS Store URL</Label>
                    <Input
                      id="ios-store-url"
                      value={iosStoreUrl}
                      onChange={(e) => setIosStoreUrl(e.target.value)}
                      placeholder="https://apps.apple.com/app/id..."
                    />
                  </div>
                </div>

                <div className="flex justify-end">
                  <Button
                    className="gap-2"
                    onClick={handleSaveSystemSettings}
                    disabled={isSavingSystem}
                  >
                    {isSavingSystem ? (
                      <Loader2 className="h-4 w-4 animate-spin" />
                    ) : (
                      <Save className="h-4 w-4" />
                    )}
                    {isSavingSystem ? 'Saving...' : 'Save Update Notice'}
                  </Button>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Database className="h-5 w-5" />
                  Legacy Deleted Cleanup
                </CardTitle>
                <CardDescription>
                  Use this tool to purge old soft-deleted/zombie records from the database.
                  Run a dry-run first, then apply.
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex flex-wrap items-center gap-2">
                  <Button
                    variant="outline"
                    onClick={() => handleOpenPurgeDialog(true)}
                    disabled={isPurging}
                  >
                    {isPurging && purgeDryRun ? (
                      <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                    ) : (
                      <Database className="mr-2 h-4 w-4" />
                    )}
                    Run Dry-Run
                  </Button>
                  <Button
                    variant="destructive"
                    onClick={() => handleOpenPurgeDialog(false)}
                    disabled={isPurging}
                  >
                    {isPurging && !purgeDryRun ? (
                      <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                    ) : (
                      <Trash2 className="mr-2 h-4 w-4" />
                    )}
                    Apply Cleanup
                  </Button>
                </div>

                {purgeResult && (
                  <div className="rounded-lg border bg-muted/30 p-4 text-sm">
                    <p className="mb-2 font-medium">
                      Last result ({purgeResult.dry_run ? 'Dry-Run' : 'Applied'}):
                    </p>
                    <div className="grid gap-1 md:grid-cols-2">
                      <p>Zombie users found: {purgeResult.zombie_users_found}</p>
                      <p>Zombie users deleted: {purgeResult.zombie_users_deleted}</p>
                      <p>Vehicle unbound count: {purgeResult.vehicle_unbound_count}</p>
                      <p>FCM tokens deleted: {purgeResult.fcm_tokens_deleted_count}</p>
                      <p>Vehicles deleted: {purgeResult.vehicles_deleted_count}</p>
                      <p>Price rules deleted: {purgeResult.price_rules_deleted}</p>
                      <p>Promotions deleted: {purgeResult.promotions_deleted}</p>
                      <p>Ratings deleted: {purgeResult.ratings_deleted}</p>
                    </div>
                  </div>
                )}
              </CardContent>
            </Card>
          </div>

          {/* Confirmation Dialog for enabling maintenance mode */}
          <AlertDialog open={showMaintenanceConfirm} onOpenChange={setShowMaintenanceConfirm}>
            <AlertDialogContent>
              <AlertDialogHeader>
                <AlertDialogTitle className="flex items-center gap-2 text-red-600">
                  <AlertTriangle className="h-5 w-5" />
                  Enable Maintenance Mode?
                </AlertDialogTitle>
                <AlertDialogDescription className="space-y-3">
                  <p>
                    This will immediately block all mobile app users from accessing the service.
                    They will see a maintenance screen instead.
                  </p>
                  <p className="font-medium text-foreground">
                    Active rides (in-progress) will NOT be interrupted.
                  </p>
                  <div className="space-y-2 pt-2">
                    <Label htmlFor="confirm-text" className="text-sm font-medium">
                      Type <span className="font-mono text-red-600">MAINTENANCE</span> to confirm:
                    </Label>
                    <Input
                      id="confirm-text"
                      value={confirmText}
                      onChange={(e) => setConfirmText(e.target.value)}
                      placeholder="MAINTENANCE"
                      className="font-mono"
                    />
                  </div>
                </AlertDialogDescription>
              </AlertDialogHeader>
              <AlertDialogFooter>
                <AlertDialogCancel onClick={() => setConfirmText('')}>Cancel</AlertDialogCancel>
                <AlertDialogAction
                  onClick={handleConfirmMaintenance}
                  disabled={confirmText !== 'MAINTENANCE'}
                  className="bg-red-600 hover:bg-red-700 disabled:opacity-50"
                >
                  Enable Maintenance Mode
                </AlertDialogAction>
              </AlertDialogFooter>
            </AlertDialogContent>
          </AlertDialog>

          <AlertDialog open={showPurgeConfirm} onOpenChange={setShowPurgeConfirm}>
            <AlertDialogContent>
              <AlertDialogHeader>
                <AlertDialogTitle className="flex items-center gap-2 text-red-600">
                  <AlertTriangle className="h-5 w-5" />
                  {purgeDryRun ? 'Run Cleanup Dry-Run?' : 'Apply Permanent Cleanup?'}
                </AlertDialogTitle>
                <AlertDialogDescription className="space-y-3">
                  {purgeDryRun ? (
                    <p>
                      This will only preview how many records would be cleaned.
                      No data will be deleted.
                    </p>
                  ) : (
                    <p>
                      This will permanently delete legacy soft-deleted/zombie records.
                      This action cannot be undone.
                    </p>
                  )}
                  <div className="space-y-2 pt-2">
                    <Label htmlFor="purge-confirm-text" className="text-sm font-medium">
                      Type <span className="font-mono text-red-600">PURGE_LEGACY_DELETED</span> to confirm:
                    </Label>
                    <Input
                      id="purge-confirm-text"
                      value={purgeConfirmText}
                      onChange={(e) => setPurgeConfirmText(e.target.value)}
                      placeholder="PURGE_LEGACY_DELETED"
                      className="font-mono"
                    />
                  </div>
                </AlertDialogDescription>
              </AlertDialogHeader>
              <AlertDialogFooter>
                <AlertDialogCancel onClick={() => setPurgeConfirmText('')}>Cancel</AlertDialogCancel>
                <AlertDialogAction
                  onClick={handleRunPurgeCleanup}
                  disabled={purgeConfirmText !== 'PURGE_LEGACY_DELETED' || isPurging}
                  className="bg-red-600 hover:bg-red-700 disabled:opacity-50"
                >
                  {isPurging ? (
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  ) : null}
                  {purgeDryRun ? 'Run Dry-Run' : 'Apply Cleanup'}
                </AlertDialogAction>
              </AlertDialogFooter>
            </AlertDialogContent>
          </AlertDialog>
        </TabsContent>
      </Tabs>
    </div>
  );
}
