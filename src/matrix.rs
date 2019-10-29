struct Mat4 {
    data [f32; 16],
}

impl Mat4 {
    fn new(data: [f32; 16]) -> Mat4 {
        Mat4 { data: data }
    }
}