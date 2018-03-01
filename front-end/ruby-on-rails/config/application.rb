require_relative 'boot'

require 'csv'
require 'trie'
require 'rails/all'

# Require the gems listed in Gemfile, including any gems
# you've limited to :test, :development, or :production.
Bundler.require(*Rails.groups)


module Pharmacodb
  class Application < Rails::Application
    # Settings in config/environments/* take precedence over those specified here.
    # Application configuration should go into files in config/initializers
    # -- all .rb files in that directory are automatically loaded.

    Raven.configure do |config|
      config.dsn = 'https://af3b5f3dfe7f42908b67ad22e7697f16:8c6963095f624e43be644f5193d1ab02@sentry.io/288334'
    end

    config.exceptions_app = self.routes
    config.assets.paths << Rails.root.join("app", "assets", "files")
    config.assets.paths << Rails.root.join("app", "assets", "fonts")
    config.my_search_split = /(?i)\s+(and)?\s*|\s+with?\s*|\s+plus?\s*|\s*[&\+]+\s*|\s+/

  end
end


