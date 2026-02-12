import { z } from 'zod';

export const governmentFundingSchema = z.object({
  name: z.string().min(1).max(255),
  state: z.string().min(1),
});

export const governmentFundingPeriodSchema = z
  .object({
    from: z.string().min(1),
    to: z.string().optional(),
    comment: z.string().optional(),
  })
  .refine((data) => !data.to || data.to >= data.from, {
    path: ['to'],
    message: 'End date must be after start date',
  });

export const governmentFundingPropertySchema = z.object({
  key: z.string().min(1),
  value: z.string().min(1),
  payment: z.number().min(0),
  requirement: z.number().min(0),
  min_age: z.number().optional().nullable(),
  max_age: z.number().optional().nullable(),
  comment: z.string().optional(),
});

export type GovernmentFundingFormData = z.infer<typeof governmentFundingSchema>;
export type GovernmentFundingPeriodFormData = z.infer<typeof governmentFundingPeriodSchema>;
export type GovernmentFundingPropertyFormData = z.infer<typeof governmentFundingPropertySchema>;
