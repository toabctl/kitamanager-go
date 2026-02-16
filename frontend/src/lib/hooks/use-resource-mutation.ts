import { useMutation, useQueryClient, type QueryKey } from '@tanstack/react-query';
import { useTranslations } from 'next-intl';
import { useToast } from '@/lib/hooks/use-toast';
import { getErrorMessage } from '@/lib/api/client';

interface ResourceMutationConfig<TData, TResponse = unknown> {
  /** The mutation function to call. */
  mutationFn: (data: TData) => Promise<TResponse>;
  /** Query key to invalidate on success. */
  invalidateQueryKey: QueryKey;
  /** Toast message shown on success. */
  successMessage: string;
  /** Fallback error message shown on failure. */
  errorMessage: string;
  /** Called after a successful mutation (e.g., close dialog, reset form). */
  onSuccess?: () => void;
}

/**
 * Lightweight mutation hook for nested resource operations on detail pages.
 * Wraps useMutation with automatic query invalidation, success/error toasts,
 * and an optional onSuccess callback for UI cleanup.
 */
export function useResourceMutation<TData, TResponse = unknown>(
  config: ResourceMutationConfig<TData, TResponse>
) {
  const queryClient = useQueryClient();
  const t = useTranslations();
  const { toast } = useToast();

  return useMutation({
    mutationFn: config.mutationFn,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: config.invalidateQueryKey });
      toast({ title: config.successMessage });
      config.onSuccess?.();
    },
    onError: (error: unknown) => {
      toast({
        title: t('common.error'),
        description: getErrorMessage(error, config.errorMessage),
        variant: 'destructive',
      });
    },
  });
}
