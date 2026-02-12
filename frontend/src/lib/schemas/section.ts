import { z } from 'zod';

export const sectionSchema = z.object({
  name: z.string().min(1).max(255),
});

export type SectionFormData = z.infer<typeof sectionSchema>;
