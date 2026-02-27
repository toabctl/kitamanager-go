import { toast } from '@/lib/hooks/use-toast';
import { getErrorMessage } from '@/lib/api/client';

/**
 * Show a destructive error toast with the API error message extracted,
 * or falling back to the given message.
 *
 * Uses the module-level toast() function, so it can be called outside
 * of React components (e.g., in mutation onError callbacks).
 */
export function showErrorToast(title: string, error: unknown, fallbackMessage: string) {
  toast({
    title,
    description: getErrorMessage(error, fallbackMessage),
    variant: 'destructive',
  });
}
