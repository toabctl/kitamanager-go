import { z } from 'zod';

export const personBaseSchema = z.object({
  first_name: z.string().min(1).max(255),
  last_name: z.string().min(1).max(255),
  gender: z.enum(['male', 'female', 'diverse']),
  birthdate: z.string().min(1),
});

export type PersonFormData = z.infer<typeof personBaseSchema>;
