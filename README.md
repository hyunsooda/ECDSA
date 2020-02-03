# ECDSA
ECDSA Keypair Generator

<ul>
 <li> Signing can be also implemented same as keypair generator mechanism (P = e*G) or just choose within order number  </li>
 <li> Verification can be also implemented using simple formula which takes private key(e) and message hash(z) within 5 lines! </li>
 <li> K, cryptographically random value, is used to both generate keypair and sign.  <code> K := max.Exp(big2, big.NewInt(k), n).Sub(max, big.NewInt(1)) </code> // (2^k % n) -1, 
 where k is between 0 and 256 and n = order parameter as according to the secp256k1 specification (see: https://en.bitcoin.it/wiki/Secp256k1) </li>
</ul>
