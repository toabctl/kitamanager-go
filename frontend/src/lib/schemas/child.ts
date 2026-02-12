import { z } from 'zod';

export const childSchema = z.object({
  first_name: z.string().min(1),
  last_name: z.string().min(1),
  gender: z.enum(['male', 'female', 'diverse']),
  birthdate: z.string().min(1),
});

export const childContractSchema = z
  .object({
    from: z.string().min(1),
    to: z.string().optional(),
    properties: z.record(z.string()).optional(),
  })
  .refine((data) => !data.to || data.to >= data.from, {
    path: ['to'],
    message: 'End date must be after start date',
  });

export type ChildFormData = z.infer<typeof childSchema>;
export type ChildContractFormData = z.infer<typeof childContractSchema>;
