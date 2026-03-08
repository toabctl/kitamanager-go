import { userCreateSchema, userUpdateSchema } from '../user';

describe('userCreateSchema', () => {
  const valid = {
    name: 'John Doe',
    email: 'john@example.com',
    password: 'secret123',
    active: true,
  };

  it('accepts valid user data', () => {
    expect(userCreateSchema.safeParse(valid).success).toBe(true);
  });

  it('rejects empty name', () => {
    expect(userCreateSchema.safeParse({ ...valid, name: '' }).success).toBe(false);
  });

  it('rejects invalid email', () => {
    expect(userCreateSchema.safeParse({ ...valid, email: 'not-an-email' }).success).toBe(false);
  });

  it('rejects password shorter than 8 chars', () => {
    expect(userCreateSchema.safeParse({ ...valid, password: 'short' }).success).toBe(false);
  });

  it('rejects password longer than 72 chars', () => {
    expect(userCreateSchema.safeParse({ ...valid, password: 'A'.repeat(73) }).success).toBe(false);
  });

  it('accepts 8-char password', () => {
    expect(userCreateSchema.safeParse({ ...valid, password: '12345678' }).success).toBe(true);
  });

  it('defaults active to true', () => {
    const { active, ...rest } = valid;
    const result = userCreateSchema.parse(rest);
    expect(result.active).toBe(true);
  });
});

describe('userUpdateSchema', () => {
  it('accepts valid update', () => {
    expect(
      userUpdateSchema.safeParse({ name: 'Jane', email: 'jane@example.com', active: false }).success
    ).toBe(true);
  });

  it('rejects invalid email', () => {
    expect(userUpdateSchema.safeParse({ name: 'Jane', email: 'bad', active: true }).success).toBe(
      false
    );
  });

  it('rejects empty name', () => {
    expect(
      userUpdateSchema.safeParse({ name: '', email: 'jane@example.com', active: true }).success
    ).toBe(false);
  });
});
