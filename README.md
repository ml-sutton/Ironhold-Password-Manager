# Ironhold-Password-Manager
A zero-knowledge, self-hosted password manager with client-side encryption. Own your credentials, own your data.

--- 

## Salts
Salts are randomly generated values combined with the master password before being passed into Argon2id. Their purpose is uniqueness rather than secrecy. Ironhold uses two distinct salts to ensure the encryption key and authentication key are derived independently from the same master password. Both are 32 bytes and stored server-side.

### SALT_A 

SALT_A is generated once at account registration and is consistent across all of a user's devices. It is combined with the master password and passed into Argon2id to produce the 32-byte AES-256 encryption key. Being stored server-side means any authenticated device can retrieve it and derive the same key.

### SALT_B

SALT_B is generated at device registration and is unique to each device. It is combined with the master password and passed into Argon2id to produce the seed for that device's ed25519 keypair. Each device therefore holds a distinct keypair, and revoking a device means deleting its public key server-side.

## Authentication 

When deciding on an authentication strategy we turned to SSH. One algorithm stood out to us; ed25519. Other strategies were evaluated such as JSON Web Tokens (JWTs) and Passkeys and WebAuthn. We decided against JWTs due to how JWTs are minted. A user must send their password over the wire via HTTP/HTTPS. We saw this as opening up an attack surface. Passkeys are being considered as a future addition when the base system is working as intended. ed25519 presented itself as a strong fit for our use case. By using a master password and a secondary salt `SALT_B` we are able to generate asymmetric keys.

### Authentication Flow


### Documentation
- [Understanding ed25519 Messari.io](https://messari.io/copilot/share/understanding-ed25519-ee5c6e19-3f05-4557-ac8f-9a1e84b9e8ff)
- [Reference implementation of ed25519 Eiken.dev](https://www.eiken.dev/blog/2020/11/code-spotlight-the-reference-implementation-of-ed25519-part-1/)
--- 

## Encryption 
For encryption we use AES-256-GCM. Other encryption algorithms were considered, a good example is XChaCha20-Poly1305. We chose AES-256-GCM as it is natively supported inside of web crypto saving us a dependency. AES-256-GCM has a smaller nonce space of 96 bits in comparison to XChaCha20-Poly1305's 192 bit nonce space. The smaller 96 bit nonce space does mean that there is more of a chance for generating duplicate nonces but a 96 bit nonce space has a birthday bound of 2^48 encryptions which for our use case is negligible. Which led us to favour AES-256-GCM. Below you will find our encryption and decryption flow which would run in the browser extension.




### Encryption Flow 

#### Steps
1. User enters their master password into the browser extension 
2. The browser extension combines the master password with `SALT_A` and runs it through a key derivation function (KDF). In our case Argon2id. This produces a 32 byte AES-256 Key which never leaves the browser extension.
3. A fresh random 12 byte nonce is generated. 
4. The browser extension generates a Universal Unique Identifier (UUID) for the entry, as this UUID is unique to each entry we will call this `entry_id`. This `entry_id` and the user's `user_id` are set as Additional Authenticated Data (AAD)
5. The browser extension using WebCrypto uses AES-256-GCM to encrypt the plaintext vault entry using the key, nonce and AAD. This produces ciphertext with a 16 bytes auth tag.
6. The final blob which will be stored in the database is assembled in the format of `nonce (12 bytes) + ciphertext + auth tag (16 bytes)`
7. The encrypted blob is sent to the Golang backend. The backend REST API has no knowledge of the blob's contents. It just receives the encrypted blob and its `entry_id`.

### Decryption Flow 
Our decryption flow assumes that the vault is open and the master password has already been input into the extension.
#### Steps
1. The blob is split into its three components: the nonce (first 12 bytes), the ciphertext and the auth tag (last 16 bytes).
2. The browser extension combines the master password with `SALT_A` and runs it through a key derivation function (KDF). In our case Argon2id. This produces a 32 byte AES-256 Key which never leaves the browser extension.
3. The AAD is reconstructed from the entry's UUID (entry_id) and the user's UUID (user_id). It must match what was used during encryption exactly.
4. The browser extension using WebCrypto uses AES-256-GCM and attempts to decrypt using the key, nonce, ciphertext, Additional Authenticated Data, and auth tag.
5. The auth tag is recomputed by GCM internally and is compared against the auth tag stored in the blob.
6. If the auth tags match, decryption succeeds and the plaintext is returned to the browser extension. If they do not match, decryption is hard rejected and no plaintext is ever produced.

### Argon2id

Argon2id is the key derivation function used to derive cryptographic keys from the master password. It was chosen over alternatives such as PBKDF2 due to its memory hardness. PBKDF2 is CPU-bound, meaning an attacker with specialised hardware such as GPUs or ASICs can run many parallel attempts cheaply. Argon2id requires a configurable amount of RAM per attempt, making large-scale brute force significantly more expensive.

The following parameters are applied consistently across all key derivation operations:

```
KDF:            Argon2id
Memory (m):     65536 (64 MiB)
Iterations (t): 3
Parallelism (p): 4
Salt length:    32 bytes
Output length:  32 bytes
```

The 32-byte output produces either the AES-256 encryption key when combined with `SALT_A`, or the ed25519 keypair seed when combined with `SALT_B`.

These parameters must remain consistent. Changing them produces a different output from the same password and salt, making existing vault entries permanently unreadable.

--- 

## Password / Passphrase Generation
When looking at how we would generate secure passwords I found inspiration in TNoodle, TNoodle is a Java application maintained and used by the World Cube Association (WCA) to generate scrambles for speedcubing competitions. TNoodle generates scrambles that are unbiased and cryptographically fair, which matters in a fair competition. This inspired us to use rejection sampling to eliminate the modulo bias in password / passphrase generation. This ensures that each character in our charset, and each word in our word list have an equal probability of selection.

### Character sets

### Documentation 
- [Modulo Biases and how to avoid them! romailler.ch](https://romailler.ch/2020/07/28/crypto-modulo_bias_guide/)
- [TNoodle github.com](https://github.com/thewca/tnoodle)


## Environment Variables
DATABASE_HOST
DATABASE_PORT
DATABASE_NAME
POSTGRES_PASSWORD
API_CLIENT_USER
API_CLIENT_PASS
DASHBOARD_CLIENT_USER
DASHBOARD_CLIENT_PASSWORD


## Links 
- [Backend](./backend/README.md)
- [Browser Extensions](./extensions/README.md)
- [Scripts](./scripts/README.md)
- [Config Files](./config/README.md)