import { z } from 'zod';

export const groupSchema = z.object({
  name: z.string().min(1).max(255),
  active: z.boolean().default(true),
});

export type GroupFormData = z.infer<typeof groupSchema>;
