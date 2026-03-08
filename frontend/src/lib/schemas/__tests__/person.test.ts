import { personBaseSchema } from '../person';

describe('personBaseSchema', () => {
  const valid = {
    first_name: 'Emma',
    last_name: 'Schmidt',
    gender: 'female' as const,
    birthdate: '2020-03-10',
  };

  it('accepts valid person data', () => {
    expect(personBaseSchema.safeParse(valid).success).toBe(true);
  });

  it('rejects empty first_name', () => {
    expect(personBaseSchema.safeParse({ ...valid, first_name: '' }).success).toBe(false);
  });

  it('rejects empty last_name', () => {
    expect(personBaseSchema.safeParse({ ...valid, last_name: '' }).success).toBe(false);
  });

  it('rejects invalid gender', () => {
    expect(personBaseSchema.safeParse({ ...valid, gender: 'other' }).success).toBe(false);
  });

  it('accepts all valid genders', () => {
    for (const gender of ['male', 'female', 'diverse']) {
      expect(personBaseSchema.safeParse({ ...valid, gender }).success).toBe(true);
    }
  });

  it('rejects empty birthdate', () => {
    expect(personBaseSchema.safeParse({ ...valid, birthdate: '' }).success).toBe(false);
  });

  it('rejects missing fields', () => {
    expect(personBaseSchema.safeParse({}).success).toBe(false);
  });

  it('rejects first_name over 255 chars', () => {
    expect(personBaseSchema.safeParse({ ...valid, first_name: 'A'.repeat(256) }).success).toBe(
      false
    );
  });
});
