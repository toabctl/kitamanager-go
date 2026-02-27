import { useMutation, useQueryClient, type QueryKey } from '@tanstack/react-query';
import { useTranslations } from 'next-intl';
import { useToast } from '@/lib/hooks/use-toast';
import { showErrorToast } from '@/lib/utils/show-error-toast';
import { getActiveContract } from '@/lib/utils/contracts';

interface ContractMutationConfig<TCreateData, TUpdateData, TContract> {
  /** Create a new contract for an entity. */
  createFn: (entityId: number, data: TCreateData) => Promise<TContract>;
  /** Update (amend) an existing contract, atomically closing it and creating a new one. */
  updateFn: (entityId: number, contractId: number, data: TUpdateData) => Promise<TContract>;
  /** Convert create data to update data (typically strips the 'from' field). */
  toUpdateData: (createData: TCreateData) => TUpdateData;
  /** Query keys to invalidate on success. */
  invalidateQueryKeys: QueryKey[];
  /** Additional query keys to invalidate, given the entity ID (e.g., per-entity contract list). */
  extraInvalidateKeys?: (entityId: number) => QueryKey[];
  /** Called after a successful mutation (use for closing dialogs, resetting state, etc.). */
  onSuccess?: () => void;
}

interface ContractMutationVariables<TCreateData> {
  entityId: number;
  data: TCreateData;
  /** The entity with its contracts array, used to find the active contract. */
  entity: { contracts?: Array<{ id: number; from: string; to?: string | null }> } | null;
  /** Whether to end the current contract (update) instead of creating a new one. */
  endCurrentContract: boolean;
}

/**
 * Shared hook for contract create/amend mutations.
 * Handles the "end current contract and create new one" pattern used by
 * both employee and child contract management.
 */
export function useContractMutation<TCreateData, TUpdateData, TContract>(
  config: ContractMutationConfig<TCreateData, TUpdateData, TContract>
) {
  const queryClient = useQueryClient();
  const t = useTranslations();
  const { toast } = useToast();

  return useMutation({
    mutationFn: async (variables: ContractMutationVariables<TCreateData>) => {
      const { entityId, data, entity, endCurrentContract } = variables;

      if (entity && endCurrentContract) {
        const active = getActiveContract(entity.contracts);
        if (active) {
          return config.updateFn(entityId, active.id, config.toUpdateData(data));
        }
      }
      return config.createFn(entityId, data);
    },
    onSuccess: (_data, variables) => {
      for (const key of config.invalidateQueryKeys) {
        queryClient.invalidateQueries({ queryKey: key });
      }
      if (config.extraInvalidateKeys) {
        for (const key of config.extraInvalidateKeys(variables.entityId)) {
          queryClient.invalidateQueries({ queryKey: key });
        }
      }

      toast({
        title: variables.endCurrentContract
          ? t('contracts.previousContractEnded')
          : t('contracts.createSuccess'),
      });

      config.onSuccess?.();
    },
    onError: (error: unknown) => {
      showErrorToast(
        t('common.error'),
        error,
        t('common.failedToCreate', { resource: 'contract' })
      );
    },
  });
}
