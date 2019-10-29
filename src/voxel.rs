use super::vector;

pub struct Voxel {
    pub loc: vector::Vector,
    pub col: [vector::Vector; 6],
    pub data: [f32; 216],
}

impl Voxel {
    pub fn new(loc: vector::Vector, col: [vector::Vector; 6]) -> Voxel{
        let (x, y, z) = (loc.x, loc.y, loc.z);
        let (a, b, c) = (loc.x+1.0, loc.y+1.0, loc.z-1.0);
        let (r1, g1, b1) = (col[0].x, col[0].y, col[0].z);
        let (r2, g2, b2) = (col[1].x, col[1].y, col[1].z);
        let (r3, g3, b3) = (col[2].x, col[2].y, col[2].z);
        let (r4, g4, b4) = (col[3].x, col[3].y, col[3].z);
        let (r5, g5, b5) = (col[4].x, col[4].y, col[4].z);
        let (r6, g6, b6) = (col[5].x, col[5].y, col[5].z);

        let data: [f32; 216] = [
            x, y, z,
            r1, g1, b1,
            a, y, z,
            r1, g1, b1,
            a, b, z,
            r1, g1, b1, //1
            x, y, z,
            r1, g1, b1,
            a, b, z,
            r1, g1, b1,
            x, b, z,
            r1, g1, b1, //2
            a, y, z,
            r2, g2, b2,
            a, b, c,
            r2, g2, b2,
            a, b, z,
            r2, g2, b2, //3
            a, y, z,
            r2, g2, b2,
            a, y, c,
            r2, g2, b2,
            a, b, c,
            r2, g2, b2, //4
            a, y, c,
            r3, g3, b3,
            x, b, c,
            r3, g3, b3,
            a, b, c,
            r3, g3, b3, //5
            a, y, c,
            r3, g3, b3,
            x, y, c,
            r3, g3, b3,
            x, b, c,
            r3, g3, b3, //6
            x, y, c,
            r4, g4, b4,
            x, b, z,
            r4, g4, b4,
            x, b, c,
            r4, g4, b4, //7
            x, y, c,
            r4, g4, b4,
            x, y, z,
            r4, g4, b4,
            x, b, z,
            r4, g4, b4, //8
            a, b, c,
            r5, g5, b5,
            x, b, c,
            r5, g5, b5,
            x, b, z,
            r5, g5, b5, //9
            a, b, c,
            r5, g5, b5,
            x, b, z,
            r5, g5, b5,
            a, b, z,
            r5, g5, b5, //10
            a, y, z,
            r6, g6, b6,
            x, y, c,
            r6, g6, b6,
            a, y, c,
            r6, g6, b6, //11
            a, y, z,
            r6, g6, b6,
            x, y, z, 
            r6, g6, b6,
            x, y, c,
            r6, g6, b6, //12
        ];

        Voxel { loc: loc, col: col, data: data }
    }
}