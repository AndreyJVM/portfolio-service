package ru.vorobevaqa.dto;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.Data;

@Data
@JsonIgnoreProperties(ignoreUnknown = true)
public class AnsibleRoleDto {

    private String name;
    private String description;

    @JsonProperty("download_count")
    private int downloadCount;

    @JsonProperty("github_user")
    private String githubUser;

    @JsonProperty("github_repo")
    private String githubRepo;

    public String getUrl() {
        return "https://galaxy.ansible.com/ui/standalone/roles/" + githubUser + "/" + name + "/";
    }
}
