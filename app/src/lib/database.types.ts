export type Json =
  | string
  | number
  | boolean
  | null
  | { [key: string]: Json | undefined }
  | Json[]

export type Database = {
  public: {
    Tables: {
      github_connection_states: {
        Row: {
          created_at: string | null
          expires_at: string
          id: string
          state_token: string
          user_id: string | null
        }
        Insert: {
          created_at?: string | null
          expires_at: string
          id?: string
          state_token: string
          user_id?: string | null
        }
        Update: {
          created_at?: string | null
          expires_at?: string
          id?: string
          state_token?: string
          user_id?: string | null
        }
        Relationships: [
          {
            foreignKeyName: "github_connection_states_user_id_fkey"
            columns: ["user_id"]
            isOneToOne: false
            referencedRelation: "users"
            referencedColumns: ["user_id"]
          },
        ]
      }
      github_installations: {
        Row: {
          id: string
          installation_id: number
          installed_at: string | null
          suspended: boolean | null
          user_id: string | null
        }
        Insert: {
          id?: string
          installation_id: number
          installed_at?: string | null
          suspended?: boolean | null
          user_id?: string | null
        }
        Update: {
          id?: string
          installation_id?: number
          installed_at?: string | null
          suspended?: boolean | null
          user_id?: string | null
        }
        Relationships: [
          {
            foreignKeyName: "github_installations_user_id_fkey"
            columns: ["user_id"]
            isOneToOne: false
            referencedRelation: "users"
            referencedColumns: ["user_id"]
          },
        ]
      }
      github_repositories: {
        Row: {
          auto_deploy_main: boolean | null
          created_at: string | null
          dockerfile_path: string | null
          enable_preview_deployments: boolean | null
          gcp_project_id: string | null
          id: string
          installation_id: number
          is_private: boolean | null
          repo_id: number
          repo_name: string
          service_name: string | null
          user_id: string | null
        }
        Insert: {
          auto_deploy_main?: boolean | null
          created_at?: string | null
          dockerfile_path?: string | null
          enable_preview_deployments?: boolean | null
          gcp_project_id?: string | null
          id?: string
          installation_id: number
          is_private?: boolean | null
          repo_id: number
          repo_name: string
          service_name?: string | null
          user_id?: string | null
        }
        Update: {
          auto_deploy_main?: boolean | null
          created_at?: string | null
          dockerfile_path?: string | null
          enable_preview_deployments?: boolean | null
          gcp_project_id?: string | null
          id?: string
          installation_id?: number
          is_private?: boolean | null
          repo_id?: number
          repo_name?: string
          service_name?: string | null
          user_id?: string | null
        }
        Relationships: [
          {
            foreignKeyName: "github_repositories_installation_fkey"
            columns: ["installation_id"]
            isOneToOne: false
            referencedRelation: "github_installations"
            referencedColumns: ["installation_id"]
          },
          {
            foreignKeyName: "github_repositories_user_id_fkey"
            columns: ["user_id"]
            isOneToOne: false
            referencedRelation: "users"
            referencedColumns: ["user_id"]
          },
        ]
      }
      users: {
        Row: {
          created_at: string | null
          github_connected: boolean
          google_refresh_token: string | null
          id: string
          updated_at: string | null
          user_id: string | null
        }
        Insert: {
          created_at?: string | null
          github_connected?: boolean
          google_refresh_token?: string | null
          id?: string
          updated_at?: string | null
          user_id?: string | null
        }
        Update: {
          created_at?: string | null
          github_connected?: boolean
          google_refresh_token?: string | null
          id?: string
          updated_at?: string | null
          user_id?: string | null
        }
        Relationships: []
      }
    }
    Views: {
      [_ in never]: never
    }
    Functions: {
      add_repositories: {
        Args: { p_installation_id: number; p_repositories: Json }
        Returns: undefined
      }
      delete_installation: {
        Args: { p_installation_id: number }
        Returns: undefined
      }
      handle_github_installation: {
        Args: {
          p_user_id: string
          p_installation_id: number
          p_repositories: Json
        }
        Returns: Json
      }
      remove_repositories: {
        Args: { p_installation_id: number; p_repository_names: string[] }
        Returns: undefined
      }
      resume_installation: {
        Args: { p_installation_id: number }
        Returns: undefined
      }
      suspend_installation: {
        Args: { p_installation_id: number }
        Returns: undefined
      }
    }
    Enums: {
      [_ in never]: never
    }
    CompositeTypes: {
      [_ in never]: never
    }
  }
}

type DefaultSchema = Database[Extract<keyof Database, "public">]

export type Tables<
  DefaultSchemaTableNameOrOptions extends
    | keyof (DefaultSchema["Tables"] & DefaultSchema["Views"])
    | { schema: keyof Database },
  TableName extends DefaultSchemaTableNameOrOptions extends {
    schema: keyof Database
  }
    ? keyof (Database[DefaultSchemaTableNameOrOptions["schema"]]["Tables"] &
        Database[DefaultSchemaTableNameOrOptions["schema"]]["Views"])
    : never = never,
> = DefaultSchemaTableNameOrOptions extends { schema: keyof Database }
  ? (Database[DefaultSchemaTableNameOrOptions["schema"]]["Tables"] &
      Database[DefaultSchemaTableNameOrOptions["schema"]]["Views"])[TableName] extends {
      Row: infer R
    }
    ? R
    : never
  : DefaultSchemaTableNameOrOptions extends keyof (DefaultSchema["Tables"] &
        DefaultSchema["Views"])
    ? (DefaultSchema["Tables"] &
        DefaultSchema["Views"])[DefaultSchemaTableNameOrOptions] extends {
        Row: infer R
      }
      ? R
      : never
    : never

export type TablesInsert<
  DefaultSchemaTableNameOrOptions extends
    | keyof DefaultSchema["Tables"]
    | { schema: keyof Database },
  TableName extends DefaultSchemaTableNameOrOptions extends {
    schema: keyof Database
  }
    ? keyof Database[DefaultSchemaTableNameOrOptions["schema"]]["Tables"]
    : never = never,
> = DefaultSchemaTableNameOrOptions extends { schema: keyof Database }
  ? Database[DefaultSchemaTableNameOrOptions["schema"]]["Tables"][TableName] extends {
      Insert: infer I
    }
    ? I
    : never
  : DefaultSchemaTableNameOrOptions extends keyof DefaultSchema["Tables"]
    ? DefaultSchema["Tables"][DefaultSchemaTableNameOrOptions] extends {
        Insert: infer I
      }
      ? I
      : never
    : never

export type TablesUpdate<
  DefaultSchemaTableNameOrOptions extends
    | keyof DefaultSchema["Tables"]
    | { schema: keyof Database },
  TableName extends DefaultSchemaTableNameOrOptions extends {
    schema: keyof Database
  }
    ? keyof Database[DefaultSchemaTableNameOrOptions["schema"]]["Tables"]
    : never = never,
> = DefaultSchemaTableNameOrOptions extends { schema: keyof Database }
  ? Database[DefaultSchemaTableNameOrOptions["schema"]]["Tables"][TableName] extends {
      Update: infer U
    }
    ? U
    : never
  : DefaultSchemaTableNameOrOptions extends keyof DefaultSchema["Tables"]
    ? DefaultSchema["Tables"][DefaultSchemaTableNameOrOptions] extends {
        Update: infer U
      }
      ? U
      : never
    : never

export type Enums<
  DefaultSchemaEnumNameOrOptions extends
    | keyof DefaultSchema["Enums"]
    | { schema: keyof Database },
  EnumName extends DefaultSchemaEnumNameOrOptions extends {
    schema: keyof Database
  }
    ? keyof Database[DefaultSchemaEnumNameOrOptions["schema"]]["Enums"]
    : never = never,
> = DefaultSchemaEnumNameOrOptions extends { schema: keyof Database }
  ? Database[DefaultSchemaEnumNameOrOptions["schema"]]["Enums"][EnumName]
  : DefaultSchemaEnumNameOrOptions extends keyof DefaultSchema["Enums"]
    ? DefaultSchema["Enums"][DefaultSchemaEnumNameOrOptions]
    : never

export type CompositeTypes<
  PublicCompositeTypeNameOrOptions extends
    | keyof DefaultSchema["CompositeTypes"]
    | { schema: keyof Database },
  CompositeTypeName extends PublicCompositeTypeNameOrOptions extends {
    schema: keyof Database
  }
    ? keyof Database[PublicCompositeTypeNameOrOptions["schema"]]["CompositeTypes"]
    : never = never,
> = PublicCompositeTypeNameOrOptions extends { schema: keyof Database }
  ? Database[PublicCompositeTypeNameOrOptions["schema"]]["CompositeTypes"][CompositeTypeName]
  : PublicCompositeTypeNameOrOptions extends keyof DefaultSchema["CompositeTypes"]
    ? DefaultSchema["CompositeTypes"][PublicCompositeTypeNameOrOptions]
    : never

export const Constants = {
  public: {
    Enums: {},
  },
} as const
