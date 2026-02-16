'use client';

import { useTranslations } from 'next-intl';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import type { Gender } from '@/lib/api/types';

export interface GenderSelectProps {
  value: Gender;
  onValueChange: (value: Gender) => void;
  'aria-invalid'?: boolean;
  'aria-describedby'?: string;
}

export function GenderSelect({
  value,
  onValueChange,
  'aria-invalid': ariaInvalid,
  'aria-describedby': ariaDescribedBy,
}: GenderSelectProps) {
  const t = useTranslations();

  return (
    <Select value={value} onValueChange={onValueChange}>
      <SelectTrigger aria-invalid={ariaInvalid} aria-describedby={ariaDescribedBy}>
        <SelectValue placeholder={t('gender.selectGender')} />
      </SelectTrigger>
      <SelectContent>
        <SelectItem value="male">{t('gender.male')}</SelectItem>
        <SelectItem value="female">{t('gender.female')}</SelectItem>
        <SelectItem value="diverse">{t('gender.diverse')}</SelectItem>
      </SelectContent>
    </Select>
  );
}
