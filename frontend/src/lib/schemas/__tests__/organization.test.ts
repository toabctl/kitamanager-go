import { organizationCreateSchema, organizationUpdateSchema } from '../organization';

describe('organizationCreateSchema', () => {
  const valid = {
    name: 'Kita Sonnenschein',
    state: 'berlin',
    active: true,
    default_section_name: 'Default',
  };

  it('accepts valid data', () => {
    expect(organizationCreateSchema.safeParse(valid).success).toBe(true);
  });

  it('rejects empty name', () => {
    expect(organizationCreateSchema.safeParse({ ...valid, name: '' }).success).toBe(false);
  });

  it('rejects name over 255 chars', () => {
    expect(organizationCreateSchema.safeParse({ ...valid, name: 'A'.repeat(256) }).success).toBe(
      false
    );
  });

  it('rejects empty state', () => {
    expect(organizationCreateSchema.safeParse({ ...valid, state: '' }).success).toBe(false);
  });

  it('requires default_section_name', () => {
    const { default_section_name, ...rest } = valid;
    expect(organizationCreateSchema.safeParse(rest).success).toBe(false);
  });

  it('rejects empty default_section_name', () => {
    expect(organizationCreateSchema.safeParse({ ...valid, default_section_name: '' }).success).toBe(
      false
    );
  });

  it('defaults active to true', () => {
    const { active, ...rest } = valid;
    const result = organizationCreateSchema.parse(rest);
    expect(result.active).toBe(true);
  });
});

describe('organizationUpdateSchema', () => {
  it('accepts valid update', () => {
    expect(
      organizationUpdateSchema.safeParse({ name: 'Updated', state: 'berlin', active: false })
        .success
    ).toBe(true);
  });

  it('rejects empty name', () => {
    expect(
      organizationUpdateSchema.safeParse({ name: '', state: 'berlin', active: true }).success
    ).toBe(false);
  });
});
