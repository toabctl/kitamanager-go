import { z } from 'zod';

export const organizationSchema = z.object({
  name: z.string().min(1).max(255),
  state: z.string().min(1),
  active: z.boolean().default(true),
  default_section_name: z.string().min(1).max(255).optional(),
});

export type OrganizationFormData = z.infer<typeof organizationSchema>;
