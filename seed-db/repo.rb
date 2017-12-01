require_relative './data_mixin'

class Repo
  include Virtus.model
  include DataMixin

  def record
    connection.exec("INSERT INTO repos (id, data) VALUES ($1, $2)", [id, data])
  end

  def url
    data_hash['url']
  end
end