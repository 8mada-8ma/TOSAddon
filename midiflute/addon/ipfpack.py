# ipf を作る/読む小さな道具(ToSアドオン用)
import struct, zlib
PW = bytes([0x6F,0x66,0x4F,0x31,0x61,0x30,0x75,0x65,0x58,0x41,0x3F,0x20,0x5B,0xFF,0x73,0x20,0x68,0x20,0x25,0x3F])
CRC = []
for i in range(256):
    c = i
    for _ in range(8): c = (c >> 1) ^ 0xEDB88320 if c & 1 else c >> 1
    CRC.append(c)
def _crc(c, b): return CRC[(c ^ b) & 0xff] ^ (c >> 8)
def _cipher(data, encrypt):
    k0, k1, k2 = 0x12345678, 0x23456789, 0x34567890
    def upd(b):
        nonlocal k0, k1, k2
        k0 = _crc(k0, b); k1 = ((k1 + (k0 & 0xff)) * 0x08088405 + 1) & 0xffffffff; k2 = _crc(k2, (k1 >> 24) & 0xff)
    for b in PW: upd(b)
    out = bytearray(data)
    for i in range(0, len(out), 2):
        t = (k2 | 2) & 0xffff; k = ((t * (t ^ 1)) >> 8) & 0xff
        plain = out[i] if encrypt else out[i] ^ k
        out[i] ^= k
        upd(plain)
    return bytes(out)
def decrypt(d): return _cipher(d, False)
def encrypt(d): return _cipher(d, True)

def read_entries(path):
    with open(path, 'rb') as f:
        f.seek(-24, 2); cnt, tbl, _, fp, magic, base, rev = struct.unpack('<HIHIIII', f.read(24))
        f.seek(tbl); es = []
        for _ in range(cnt):
            nl, crc, cl, ul, ptr, cnl = struct.unpack('<HIIIIH', f.read(20))
            cont = f.read(cnl).decode('utf8'); name = f.read(nl).decode('utf8')
            es.append(dict(cont=cont, name=name, crc=crc, cl=cl, ul=ul, ptr=ptr))
        for e in es:
            f.seek(e['ptr']); raw = f.read(e['cl'])
            if e['cl'] == e['ul']: e['data'] = raw; e['enc'] = False; continue
            try: e['data'] = zlib.decompress(raw, -15); e['enc'] = False
            except zlib.error: e['data'] = zlib.decompress(decrypt(raw), -15); e['enc'] = True
        return es, (magic, base, rev)

def write_ipf(path, files, enc=True):
    """files: [(container, name, bytes)]"""
    body = bytearray(); table = bytearray()
    for cont, name, data in files:
        co = zlib.compressobj(9, zlib.DEFLATED, -15); comp = co.compress(data) + co.flush()
        if enc: comp = encrypt(comp)
        ptr = len(body); body += comp
        cb, nb = cont.encode('utf8'), name.encode('utf8')
        table += struct.pack('<HIIIIH', len(nb), zlib.crc32(data), len(comp), len(data), ptr, len(cb)) + cb + nb
    tbl = len(body); fp = tbl + len(table)
    with open(path, 'wb') as f:
        f.write(body + table + struct.pack('<HIHIIII', len(files), tbl, 0, fp, 0x06054B50, 0, 0))
