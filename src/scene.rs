use wasm_bindgen::prelude::*;
use std::vec::Vec;
use web_sys::{WebGlBuffer, WebGlProgram, WebGlRenderingContext, WebGlUniformLocation};
use super::voxel;
use super::web_gl;

pub struct Scene<'a> {
    pub voxels: Vec<voxel::Voxel>,
    pub gl: &'a WebGlRenderingContext,
    pub program: WebGlProgram,
    pub transform: [f32; 16],
    u_transform_location: WebGlUniformLocation,
    a_vertex_location: u32,
    a_color_location: u32,
    vbo_id: &'a WebGlBuffer,
    vbo_len: i32,
    update: bool,
}

impl<'a> Scene<'a> {
    pub fn new(
        gl: &'a WebGlRenderingContext
        ) -> Result<Scene<'a>, JsValue> {

        //Create the program
        let vert_shader = web_gl::WebGL::compile_shader(
            &gl,
            WebGlRenderingContext::VERTEX_SHADER,
            r#"
            // VertexShader

            precision mediump int;
            precision mediump float;

            uniform mat4 u_Transform;

            attribute vec3 a_Vertex;
            attribute vec3 a_Color;

            varying vec4 v_vertex_color;

            void main() {
                gl_Position = u_Transform * vec4(a_Vertex, 1.0);

                v_vertex_color = vec4(a_Color, 1.0);
            }
            "#,
        )?;
        let frag_shader = web_gl::WebGL::compile_shader(
            &gl,
            WebGlRenderingContext::FRAGMENT_SHADER,
            r#"
            // Fragment shader

            precision mediump int;
            precision mediump float;

            varying vec4 v_vertex_color;

            void main() {
                gl_FragColor = v_vertex_color;
            }
        "#,
        )?;
        let program = web_gl::WebGL::link_program(&gl, &vert_shader, &frag_shader)?;
        gl.use_program(Some(&program));

        let ret = Scene { 
            gl: gl,
            voxels: Vec::new(),
            program: program,
            transform: [
                1.0, 0.0, 0.0, 0.0,
                0.0, 1.0, 0.0, 0.0,
                0.0, 0.0, 1.0, 0.0,
                0.0, 0.0, 0.0, 1.0,
            ],
            u_transform_location: gl.get_uniform_location(&program, "u_Transform").expect("could not get u_Transform location"),
            a_color_location: gl.get_attrib_location(&program, "a_Color") as u32,
            a_vertex_location: gl.get_attrib_location(&program,  "a_Vertex") as u32,
            update: false,
            vbo_id: &gl.create_buffer().ok_or("failed to create buffer")?,
            vbo_len: -1,
        };

        return Ok(ret);
    }

    pub fn render(&self) -> Result<(), String> {
        //Set all the uniforms
        self.gl.uniform_matrix4fv_with_f32_array(Some(&self.u_transform_location), false, &self.transform);

        //Set all the attributes
        let vbo_id = self.get_vbo()?;
        self.gl.bind_buffer(WebGlRenderingContext::ARRAY_BUFFER, Some(&vbo_id));

        let offset = 4*6;
        self.gl.vertex_attrib_pointer_with_i32(self.a_vertex_location, 3, WebGlRenderingContext::FLOAT, false, offset, 0);
        self.gl.enable_vertex_attrib_array(self.a_vertex_location);
        self.gl.vertex_attrib_pointer_with_i32(self.a_color_location, 3, WebGlRenderingContext::FLOAT, false, offset, 4*3);
        self.gl.enable_vertex_attrib_array(self.a_color_location);

        self.gl.draw_arrays(
            WebGlRenderingContext::TRIANGLES,
            0,
            self.vbo_len,
        );

        Ok(())
    }

    pub fn add_voxel(&self, vox: voxel::Voxel) {
        self.voxels.push(vox);
        self.update = true;
    }

    fn get_vbo(&self) -> Result<&WebGlBuffer, String> {
        if !self.update {
            return Ok(self.vbo_id);
        }

        let vbo_data: &[f32];
        let i = 0;
        let j = 0;
        for v in self.voxels {
            for d in v.data.iter() {
                vbo_data[i*j+i] = *d;
                j += 1;
            }
            i += 1;
        }

        self.vbo_len = vbo_data.len() as i32;

        self.update = false;
        self.vbo_id = web_gl::WebGL::update_buffer(&self.gl, &vbo_data, &self.vbo_id);
        return Ok(self.vbo_id);
    }
}