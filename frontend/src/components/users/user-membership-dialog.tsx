'use client';

import { useState } from 'react';
import { useTranslations } from 'next-intl';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Trash2 } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from '@/components/ui/dialog';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
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
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { useToast } from '@/lib/hooks/use-toast';
import { apiClient, getErrorMessage } from '@/lib/api/client';
import { queryKeys } from '@/lib/api/queryKeys';
import type { User, Role, UserMembership } from '@/lib/api/types';

const ROLES: Role[] = ['admin', 'manager', 'member'];

interface UserMembershipDialogProps {
  user: User | null;
  orgId: number;
  onClose: () => void;
}

export function UserMembershipDialog({ user, orgId, onClose }: UserMembershipDialogProps) {
  const t = useTranslations();
  const { toast } = useToast();
  const queryClient = useQueryClient();
  const [removeTarget, setRemoveTarget] = useState<UserMembership | null>(null);

  const { data: membershipsData, isLoading: membershipsLoading } = useQuery({
    queryKey: queryKeys.users.memberships(user?.id ?? 0),
    queryFn: () => apiClient.getUserMemberships(user!.id),
    enabled: !!user,
  });

  const orgMembership =
    membershipsData?.memberships?.find((m) => m.organization_id === orgId) ?? null;

  const invalidateMemberships = () => {
    if (user) {
      queryClient.invalidateQueries({ queryKey: queryKeys.users.memberships(user.id) });
    }
  };

  const updateRoleMutation = useMutation({
    mutationFn: ({ organizationId, role }: { organizationId: number; role: Role }) =>
      apiClient.updateUserOrganizationRole(user!.id, organizationId, role),
    onSuccess: () => {
      toast({ title: t('users.roleUpdated') });
      invalidateMemberships();
    },
    onError: (error) => {
      toast({
        title: t('users.failedToUpdateRole'),
        description: getErrorMessage(error, t('common.error')),
        variant: 'destructive',
      });
    },
  });

  const removeMutation = useMutation({
    mutationFn: (organizationId: number) =>
      apiClient.removeUserFromOrganization(user!.id, organizationId),
    onSuccess: () => {
      toast({ title: t('users.removedFromOrganization') });
      invalidateMemberships();
      setRemoveTarget(null);
    },
    onError: (error) => {
      toast({
        title: t('users.failedToRemoveFromOrganization'),
        description: getErrorMessage(error, t('common.error')),
        variant: 'destructive',
      });
    },
  });

  return (
    <>
      <Dialog open={!!user} onOpenChange={(open) => !open && onClose()}>
        <DialogContent className="max-w-lg">
          <DialogHeader>
            <DialogTitle>
              {user
                ? t('users.organizationMembershipFor', { name: user.name })
                : t('users.organizationMembership')}
            </DialogTitle>
            <DialogDescription className="sr-only">
              {t('users.organizationMembership')}
            </DialogDescription>
          </DialogHeader>

          {membershipsLoading ? (
            <div className="text-muted-foreground py-4 text-center text-sm">
              {t('common.loading')}
            </div>
          ) : orgMembership ? (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>{t('roles.role')}</TableHead>
                  <TableHead className="w-[60px]" />
                </TableRow>
              </TableHeader>
              <TableBody>
                <TableRow>
                  <TableCell>
                    <Select
                      value={orgMembership.role}
                      onValueChange={(role) =>
                        updateRoleMutation.mutate({
                          organizationId: orgMembership.organization_id,
                          role: role as Role,
                        })
                      }
                    >
                      <SelectTrigger className="w-[140px]">
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        {ROLES.map((role) => (
                          <SelectItem key={role} value={role}>
                            {t(`roles.${role}`)}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </TableCell>
                  <TableCell>
                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() => setRemoveTarget(orgMembership)}
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </TableCell>
                </TableRow>
              </TableBody>
            </Table>
          ) : (
            <div className="flex items-center gap-2">
              <span className="text-muted-foreground text-sm">
                {t('users.notMemberOfOrganization')}
              </span>
            </div>
          )}
        </DialogContent>
      </Dialog>

      {/* Remove confirmation */}
      <AlertDialog open={!!removeTarget} onOpenChange={(open) => !open && setRemoveTarget(null)}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>{t('users.confirmRemoval')}</AlertDialogTitle>
            <AlertDialogDescription>
              {removeTarget ? t('users.removeFromOrganizationConfirm') : ''}
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>{t('common.cancel')}</AlertDialogCancel>
            <AlertDialogAction
              onClick={() => removeTarget && removeMutation.mutate(removeTarget.organization_id)}
              disabled={removeMutation.isPending}
            >
              {t('common.delete')}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}
