# TODOs
- Remaining policies:
  - allkeyslfu
  - volatilelfu
  - allkeysrandom
  - volatilerandom
  - volatilettl
- Consider using a Policy option to select policy as the first argument to
  a single open function.
- Consider making the Store interface exported as a dependency of a
  CustomPolicy option allowing users to implement their own policies.
- Consider renaming package to "inmem" to avoid stutters.
- Update WithActiveExpiration to use an interface and make expirers public.
