import { z } from 'zod';

export const createResultSchema = z.object({
  done: z.boolean(),
  comment: z
    .string()
    .min(1, 'Комментарий обязателен')
    .max(2000, 'Комментарий не должен превышать 2000 символов'),
});

export type CreateResultFormData = z.infer<typeof createResultSchema>;

