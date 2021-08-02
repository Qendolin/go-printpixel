#version 330 core
layout (location = 0) in vec3 in_position;
layout (location = 1) in vec3 in_normal;

out vec3 pass_normal;
out vec3 pass_local_normal;
out vec3 pass_global_normal;
out vec3 pass_position;

uniform mat4 u_transform_mat;
uniform mat4 u_model_mat;

void main()
{
    gl_Position = u_transform_mat*vec4(in_position, 1.);

    pass_local_normal = in_normal;
    pass_global_normal = (u_model_mat*vec4(in_normal, 1.)).xyz;
    pass_position = gl_Position.xyz;
}